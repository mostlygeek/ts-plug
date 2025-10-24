package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "IP:port to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if required headers exist
		loginName := r.Header.Get("Tailscale-User-Login")
		displayName := r.Header.Get("Tailscale-User-Name")
		profilePicURL := r.Header.Get("Tailscale-User-Profile-Pic")

		w.Header().Set("Cache-Control", "no-cache")

		if loginName != "" && displayName != "" && profilePicURL != "" {
			// Serve HTML page with user information
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Hello from Go!</title>
</head>
<body>
    <h1>(GO) Tailscale User Information</h1>
    <p><strong>Login Name:</strong> %s</p>
    <p><strong>Name:</strong> %s</p>
    <p><strong>Profile Picture:</strong></p>
    <img src="%s" alt="Profile Picture" style="max-width: 200px;">
</body>
</html>
`, loginName, displayName, profilePicURL)
		} else {
			// Print anonymous message with IP
			ip := getClientIP(r)
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Hello anonymous from %s\n", ip)
		}
	})

	log.Printf("Starting server on %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func getClientIP(r *http.Request) string {
	// Try to get IP from X-Forwarded-For header (for proxies)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Take the first IP in the list
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Try to get IP from X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
