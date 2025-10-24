package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"tailscale.com/tsnet"
)

var (
	flagHostname   = flag.String("hostname", "tsplug", "hostname on tailnet")
	flagDir        = flag.String("dir", ".data", "directory to store tailscale state")
	flagLogLevel   = flag.String("log", "info", "Log level (debug | info | warn | error)")
	flagPort       = flag.Int("port", 8080, "port of upstream server to send traffic to")
	flagFunnel     = flag.Bool("funnel", false, "enable funnel")
	flagDebugTSNet = flag.Bool("debug-tsnet", false, "enable tsnet.Server logging")
)

func main() {
	flag.Parse()

	switch *flagLogLevel {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.Error("unknown log level", slog.String("level", *flagLogLevel))
		os.Exit(1)
	}

	// Everything after "--" goes into cmdArgs
	cmdArgs := flag.Args()
	if len(cmdArgs) == 0 {
		slog.Error("no command to run")
		return
	}

	// cmdExitChannel receives the error when cmd.Wait() return
	cmdExitChan := make(chan error)

	// signalChan receives OS signals for shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// create a context that can be cancelled to stop upstream and tsnet
	ctx, cancelCtx := context.WithCancel(context.Background())

	// start the child process that will handle requests
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = append(os.Environ(), "TSPLUG_ACTIVE=1")
	if err := attachLogging(cmd); err != nil {
		slog.Error("failed to attach logging to cmd", "error", err)
		os.Exit(1)
	}

	slog.Info("starting command", "cmd", strings.Join(cmdArgs, " "))
	if err := cmd.Start(); err != nil {
		slog.Error("command start failed", "error", err)
		os.Exit(1)
	} else {
		slog.Info("command started")
	}

	// handle the exit cases either from signal or the upstream command exiting
	go func() {
		for {
			select {
			case cmdExitChan <- cmd.Wait():
				// the upstream command has exited, we're
				return
			case sig := <-signalChan:
				slog.Info("signal received, shutting down...", "sig", sig.String())

				// this will cause the case above with cmd.Wait() to return
				// as well ts.Up() to exit early if it hasn't been fully initialized yet
				cancelCtx()
			}
		}
	}()

	ts := &tsnet.Server{
		Hostname: *flagHostname,
		Dir:      *flagDir,
	}

	if *flagDebugTSNet {
		ts.Logf = func(format string, args ...any) {
			cur := slog.SetLogLoggerLevel(slog.LevelDebug) // force debug if this option is on
			slog.Debug(fmt.Sprintf(format, args...))
			slog.SetLogLoggerLevel(cur)
		}
	}

	// start the tsnet server. important to give it a cancellable context
	// because ts only listens for SIGHUP to interrupt connecting to the tailnet
	// and causes the program to ignore SIGINT/SIGTERM. canceling the context
	// will cause ts.Up() to exit early
	st, err := ts.Up(ctx)
	if err != nil {
		slog.Error("error starting tsnet server", slog.Any("error", err))
		cancelCtx()
		os.Exit(1)
	}

	lc, err := ts.LocalClient()
	if err != nil {
		slog.Error("Failed to get tsnet LocalClient", "error", err)
		cancelCtx()
		os.Exit(1)
	}

	var tl net.Listener
	if *flagFunnel {
		tl, err = ts.ListenFunnel("tcp", ":443")
		if err != nil {
			slog.Error("failed to listen on funnel", "error", err)
			cancelCtx()
			os.Exit(1)
		}
		slog.Info(fmt.Sprintf("listening at (FUNNEL): https://%s", strings.TrimSuffix(st.Self.DNSName, ".")))
	} else {
		tl, err = ts.ListenTLS("tcp", ":443")
		if err != nil {
			slog.Error("error tailnet listening", slog.Any("error", err))
			cancelCtx()
			os.Exit(1)
		} else {
			slog.Info(fmt.Sprintf("listening at: https://%s", strings.TrimSuffix(st.Self.DNSName, ".")))
		}
	}

	u, err := url.Parse(fmt.Sprintf("http://localhost:%d", *flagPort))
	if err != nil {
		slog.Error("invalid upstream", "error", err)
		cancelCtx()
		os.Exit(1)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 2 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: time.Second,
	}

	// whoisHandler injects whois information into the request headers
	whoisHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ul, dn, pp string

		who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil {
			slog.Error("whois lookup failed", "error", err, "remote", r.RemoteAddr)
		} else if who.UserProfile != nil && who.UserProfile.LoginName != "tagged-devices" {
			slog.Debug("set Tailscale-* headers",
				slog.String("remote", r.RemoteAddr),
				slog.String("id", who.UserProfile.ID.String()),
			)

			ul = who.UserProfile.LoginName
			dn = who.UserProfile.DisplayName
			pp = who.UserProfile.ProfilePicURL
		}

		// always populate the headers, even if blank for security reasons.
		r.Header.Set("Tailscale-User-Login", ul)
		r.Header.Set("Tailscale-User-Name", dn)
		r.Header.Set("Tailscale-User-Profile-Pic", pp)

		proxy.ServeHTTP(w, r)
	})

	// Start the proxy server
	go func() {
		httpServer := &http.Server{
			Handler: whoisHandler,
		}
		httpServer.Serve(tl)
	}()

	err = <-cmdExitChan
	slog.Info("cmd exited", "error", err)
}

// attachLogging attaches logging to a command's stdout and stderr
// and logs them to the slog logger.
// It returns an error if it fails to attach the pipes.
func attachLogging(cmd *exec.Cmd) error {

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// log stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			slog.Info(fmt.Sprintf("cmd > %s", scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			slog.Error("reading stdout failed", "error", err)
		}
	}()

	// log stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			slog.Info(fmt.Sprintf("cmd stderr> %s", scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			slog.Error("reading stderr failed", "error", err)
		}
	}()

	return nil
}
