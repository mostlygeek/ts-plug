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

	// start the command
	cmdContext, cmdCancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(cmdContext, cmdArgs[0], cmdArgs[1:]...)

	// capture stdout/stderr before starting the command
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("failed to get stdout pipe", "error", err)
		os.Exit(1)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("failed to get stderr pipe", "error", err)
		os.Exit(1)
	}

	slog.Info("starting command", "cmd", strings.Join(cmdArgs, " "))
	if err := cmd.Start(); err != nil {
		slog.Error("command start failed", "error", err)
		os.Exit(1)
	} else {
		slog.Info("command started")
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
	cmdChan := make(chan struct{})
	go func() {
		cmd.Wait()
		cmdChan <- struct{}{}
	}()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)

	go func() {
		<-exitChan
		slog.Info("interrupt, shutting down upstream")
		cmdCancel()
	}()

	// start the tsnet listener
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

	st, err := ts.Up(context.Background())
	if err != nil {
		slog.Error("error starting tsnet server", slog.Any("error", err))
		os.Exit(1)
	}

	lc, err := ts.LocalClient()
	if err != nil {
		slog.Error("Failed to get tsnet LocalClient", "error", err)
		os.Exit(1)
	}

	var tl net.Listener
	if *flagFunnel {
		tl, err = ts.ListenFunnel("tcp", ":443")
		if err != nil {
			slog.Error("failed to listen on funnel", "error", err)
			os.Exit(1)
		}
		slog.Info(fmt.Sprintf("listening at (FUNNEL): https://%s", strings.TrimSuffix(st.Self.DNSName, ".")))
	} else {
		tl, err = ts.ListenTLS("tcp", ":443")
		if err != nil {
			slog.Error("error tailnet listening", slog.Any("error", err))
			os.Exit(1)
		} else {
			slog.Info(fmt.Sprintf("listening at: https://%s", strings.TrimSuffix(st.Self.DNSName, ".")))
		}
	}

	u, err := url.Parse(fmt.Sprintf("http://localhost:%d", *flagPort))
	if err != nil {
		slog.Error("invalid upstream", "error", err)
		os.Exit(1)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 2 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: time.Second,
	}

	// Create a wrapper handler that processes requests before the proxy
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client info from tailscale
		who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil {
			slog.Error("whois lookup failed", "error", err, "remote", r.RemoteAddr)
		} else if who.UserProfile != nil && who.UserProfile.LoginName != "tagged-devices" {
			slog.Debug("set Tailscale-* headers",
				slog.String("remote", r.RemoteAddr),
				slog.String("id", who.UserProfile.ID.String()),
			)
			r.Header.Set("Tailscale-User-Login", who.UserProfile.LoginName)
			r.Header.Set("Tailscale-User-Name", who.UserProfile.DisplayName)
			r.Header.Set("Tailscale-User-Profile-Pic", who.UserProfile.ProfilePicURL)
		}

		// Now pass the request to the proxy
		proxy.ServeHTTP(w, r)
	})

	go func(l net.Listener) {
		httpServer := &http.Server{
			Handler: handler,
		}
		httpServer.Serve(l)
	}(tl)

	<-cmdChan
	slog.Info("cmd exited")
}
