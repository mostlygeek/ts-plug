package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"tailscale.com/tsnet"
)

var (
	flagDir        = flag.String("dir", "", "tsnet server directory")
	flagHostname   = flag.String("hostname", "tsunplug", "hostname for the tsnet server")
	flagDebugTSNet = flag.Bool("debug-tsnet", false, "enable tsnet.Server logging")
	flagPort       = flag.Int("port", 80, "local port to listen on")
)

func main() {

	flag.Parse()
	if *flagDir == "" {
		slog.Error("dir is required")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 1 {
		slog.Error("remote-addr is required as first positional argument")
		os.Exit(1)
	}

	remoteAddr := args[0]

	// Ensure remoteAddr has a port, default to 80 if not specified
	if _, _, err := net.SplitHostPort(remoteAddr); err != nil {
		remoteAddr = net.JoinHostPort(remoteAddr, "80")
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

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

	st, err := ts.Up(ctx)
	if err != nil {
		slog.Error("error starting tsnet server", slog.Any("error", err))
		cancelCtx()
		os.Exit(1)
	}

	slog.Info("tsnet server started", slog.String("status", st.BackendState))

	target, err := url.Parse("http://" + remoteAddr)
	if err != nil {
		slog.Error("invalid remote address", slog.Any("error", err))
		os.Exit(1)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return ts.Dial(ctx, network, remoteAddr)
		},
	}

	listenAddr := fmt.Sprintf("localhost:%d", *flagPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		slog.Error("failed to listen", slog.String("addr", listenAddr), slog.Any("error", err))
		os.Exit(1)
	}
	defer listener.Close()

	slog.Info("HTTP proxy listening", slog.String("local", listenAddr), slog.String("remote", remoteAddr))

	if err := http.Serve(listener, proxy); err != nil {
		log.Fatal(err)
	}
}
