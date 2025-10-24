#!/usr/bin/env python3

import argparse
from http.server import HTTPServer, BaseHTTPRequestHandler


class TailscaleHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        # Check if required headers exist
        login_name = self.headers.get("Tailscale-User-Login")
        display_name = self.headers.get("Tailscale-User-Name")
        profile_pic_url = self.headers.get("Tailscale-User-Profile-Pic")

        self.send_response(200)
        self.send_header("Cache-Control", "no-cache")

        if login_name and display_name and profile_pic_url:
            # Serve HTML page with user information
            self.send_header("Content-Type", "text/html")
            self.end_headers()
            html = f"""
<!DOCTYPE html>
<html>
<head>
    <title>Hello from Python!</title>
</head>
<body>
    <h1>(Python) Tailscale User Information</h1>
    <p><strong>Login Name:</strong> {login_name}</p>
    <p><strong>Name:</strong> {display_name}</p>
    <p><strong>Profile Picture:</strong></p>
    <img src="{profile_pic_url}" alt="Profile Picture" style="max-width: 200px;">
</body>
</html>
"""
            self.wfile.write(html.encode())
        else:
            # Print anonymous message with IP
            ip = self.get_client_ip()
            self.send_header("Content-Type", "text/plain")
            self.end_headers()
            self.wfile.write(f"Hello anonymous from {ip}\n".encode())

    def get_client_ip(self):
        # Try to get IP from X-Forwarded-For header (for proxies)
        forwarded = self.headers.get("X-Forwarded-For")
        if forwarded:
            ips = forwarded.split(",")
            if ips:
                return ips[0].strip()

        # Try to get IP from X-Real-IP header
        real_ip = self.headers.get("X-Real-IP")
        if real_ip:
            return real_ip

        # Fall back to remote address
        return self.client_address[0]

    def log_message(self, format, *args):
        # Suppress default request logging
        pass


def main():
    parser = argparse.ArgumentParser(description="Tailscale Hello World Server")
    parser.add_argument(
        "--listen",
        type=str,
        default="localhost:8080",
        help="IP:port to listen on (default: localhost:8080)",
    )
    args = parser.parse_args()

    # Parse listen address
    listen_addr = args.listen
    if listen_addr.startswith(":"):
        host = ""
        port = int(listen_addr[1:])
    elif ":" in listen_addr:
        host, port_str = listen_addr.rsplit(":", 1)
        port = int(port_str)
    else:
        host = ""
        port = int(listen_addr)

    server_address = (host, port)
    httpd = HTTPServer(server_address, TailscaleHandler)

    display_host = host if host else "0.0.0.0"
    print(f"Starting server on {display_host}:{port}")

    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down server...")
        httpd.shutdown()


if __name__ == "__main__":
    main()
