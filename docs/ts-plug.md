# ts-plug Guide

**ts-plug** wraps your local server and exposes it on your tailnet with automatic TLS and DNS.

## Overview

ts-plug is a reverse proxy that:
- Starts and manages your upstream server process
- Connects to your tailnet
- Provides automatic HTTPS with valid TLS certificates
- Optionally exposes services publicly via Tailscale Funnel
- Supports HTTP, HTTPS, and DNS proxying

## Installation

Build from source:
```sh
make ts-plug
```

Install to $GOPATH/bin:
```sh
make install
```

## Basic Usage

The basic pattern is:
```sh
ts-plug [flags] -- [your-server-command]
```

Everything after `--` is treated as the command to run.

### Simple Examples

Run a Python HTTP server on your tailnet:
```sh
ts-plug -hostname myserver -- python -m http.server 8080
```

Run a Node.js app:
```sh
ts-plug -hostname api -- node server.js
```

Run a Go server:
```sh
ts-plug -hostname webapp -- go run main.go
```

## How It Works

```
┌─────────────────────────────────────────────────────────┐
│  Your Local Machine                                     │
│                                                         │
│  ┌──────────────┐          ┌──────────────┐             │
│  │  ts-plug     │  starts  │ Your Server  │             │
│  │              │ ──────>  │ localhost:80 │             │
│  └──────┬───────┘          └──────────────┘             │
│         │                                               │
└─────────┼───────────────────────────────────────────────┘
          │ Tailscale (encrypted)
          │
┌─────────┼───────────────────────────────────────────────┐
│  Your tailnet                                           │
│         │                                               │
│         │    ┌──────────────┐    ┌──────────────┐       │
│         └──> │   HTTPS:443  │───>│ Team Members │       │
│              │   (with TLS) │    │   Devices    │       │
│              └──────────────┘    └──────────────┘       │
└─────────────────────────────────────────────────────────┘
```

ts-plug:
1. Starts your upstream server process
2. Connects to your tailnet
3. Provisions TLS certificates automatically
4. Listens for connections on your tailnet
5. Reverse proxies traffic to your local server

## Configuration Flags

### Required

- Command after `--` - The server command to execute

### Network

- `-hostname` / `-hn` - Hostname on your tailnet (default: "tsmultiplug")
  ```sh
  ts-plug -hostname myapp -- python app.py
  # Access at: https://myapp.tailnet-name.ts.net
  ```

- `-dir` - Directory to store Tailscale state (default: ".data")
  ```sh
  ts-plug -dir /var/lib/tsplug -hostname api -- ./server
  ```

### Listeners

By default, ts-plug enables HTTPS on port 443 proxying to localhost:8080.

#### HTTP

- `-http` - Enable HTTP listener (default port mapping: 80:8080)
  ```sh
  # Enable HTTP, proxy port 80 to localhost:8080
  ts-plug -http -hostname web -- python -m http.server 8080
  ```

- `-http-port` - Customize HTTP port mapping
  ```sh
  # Listen on port 8000, proxy to localhost:3000
  ts-plug -http-port 8000:3000 -hostname web -- node server.js

  # Listen and proxy both on port 9000
  ts-plug -http-port 9000 -hostname web -- ./server
  ```

#### HTTPS

- `-https` - Enable HTTPS listener (default port mapping: 443:8080)
  ```sh
  ts-plug -https -hostname secure -- python -m http.server 8080
  ```

- `-https-port` - Customize HTTPS port mapping
  ```sh
  # Listen on port 8443, proxy to localhost:3000
  ts-plug -https-port 8443:3000 -hostname web -- node server.js
  ```

#### DNS

- `-dns` - Enable DNS listener (default port mapping: 53:53)
  ```sh
  ts-plug -dns -hostname dns -- pihole-FTL
  ```

- `-dns-port` - Customize DNS port mapping
  ```sh
  # Forward DNS from port 53 to localhost:5353
  ts-plug -dns-port 53:5353 -hostname resolver -- dnsmasq
  ```

### Public Access

- `-public` - Enable Tailscale Funnel for public HTTPS access
  ```sh
  ts-plug -public -hostname demo -- python -m http.server 8080
  # Now accessible from the public internet!
  ```

  This is perfect for:
  - Webhook testing
  - Demo sites
  - Temporary public APIs
  - Sharing work with clients

### Debugging

- `-log` - Set log level (debug, info, warn, error)
  ```sh
  ts-plug -log debug -hostname myapp -- node server.js
  ```

- `-debug-tsnet` - Enable verbose tsnet.Server logging
  ```sh
  ts-plug -debug-tsnet -hostname myapp -- ./server
  ```

## Advanced Usage

### Multiple Listeners

Enable multiple protocols simultaneously:
```sh
# HTTP, HTTPS, and DNS
ts-plug -http -https -dns -hostname multi -- ./server
```

### Custom Port Mappings

Map different ports for each protocol:
```sh
ts-plug \
  -http-port 80:3000 \
  -https-port 443:3000 \
  -hostname myapp \
  -- node server.js
```

### Environment Detection

Your server can detect when it's running under ts-plug:
```sh
if [ "$TSPLUG_ACTIVE" = "1" ]; then
  echo "Running behind ts-plug!"
fi
```

```python
import os
if os.getenv('TSPLUG_ACTIVE') == '1':
    print("Running behind ts-plug!")
```

## Security Considerations

### Automatic TLS

ts-plug automatically provisions valid TLS certificates for your tailnet hostname. No configuration needed.

### User Identity Headers

When requests come through ts-plug, these headers are added:
- `Tailscale-User-Login` - User's login email
- `Tailscale-User-Name` - User's display name
- `Tailscale-User-Profile-Pic` - URL to user's profile picture

Your server can use these for authentication:

```python
@app.route('/api/whoami')
def whoami():
    return {
        'login': request.headers.get('Tailscale-User-Login'),
        'name': request.headers.get('Tailscale-User-Name'),
        'picture': request.headers.get('Tailscale-User-Profile-Pic')
    }
```

### Network Isolation

Services are only accessible to devices on your tailnet (unless `-public` is used).

## Use Cases

### Local Development Sharing

Share your dev server with teammates:
```sh
ts-plug -hostname dev-alice -- npm run dev
# Tell your teammate to visit: https://dev-alice.tailnet.ts.net
```

### Webhook Testing

Test webhooks without ngrok:
```sh
ts-plug -public -hostname webhook-test -- python webhook_server.py
# Use the public URL in GitHub/Stripe/etc webhook settings
```

### Container Deployment

Use as an entrypoint to eliminate sidecar containers:
```dockerfile
COPY ts-plug /usr/local/bin/
ENTRYPOINT ["ts-plug", "-hostname", "myapp", "--"]
CMD ["python", "app.py"]
```

See [docker.md](./docker.md) for detailed examples.

### Multi-Protocol Services

Run Pi-hole with both DNS and HTTP:
```sh
ts-plug \
  -dns \
  -http \
  -hostname pihole \
  -- pihole-FTL
```

## Troubleshooting

### Port Already in Use

If you get "address already in use", another process is listening on the configured port:
```sh
# Check what's using port 8080
lsof -i :8080

# Use a different port
ts-plug -https-port 443:3000 -hostname myapp -- node server.js
```

### Connection Refused

If ts-plug can't connect to your server:
1. Verify your server is listening on the correct port
2. Make sure it's listening on `0.0.0.0` or `127.0.0.1`, not a specific IP
3. Check logs with `-log debug`

### Tailscale Authentication

First run will prompt you to authenticate with Tailscale:
```sh
ts-plug -hostname test -- python -m http.server 8080
# Follow the URL to authenticate
```

State is saved in the `-dir` location (default: `.data/`)

## Examples

### Next.js Development
```sh
ts-plug -hostname nextjs-dev -https-port 443:3000 -- npm run dev
```

### Django Application
```sh
ts-plug -hostname django -https-port 443:8000 -- python manage.py runserver
```

### Static Site
```sh
ts-plug -public -hostname my-site -- python -m http.server 8080
```

### API with Custom Domain
```sh
ts-plug -hostname api-v1 -https-port 443:5000 -- flask run
```

## Comparison with ts-unplug

| Feature | ts-plug | ts-unplug |
|---------|---------|-----------|
| Direction | Local → tailnet | tailnet → Local |
| Use Case | Share local services | Access remote services |
| Starts Process | Yes | No |
| TLS | Automatic | Proxies existing |
| Public Access | Optional | No |

## See Also

- [ts-unplug Guide](./ts-unplug.md) - Access remote services locally
- [Use Cases](./use-cases.md) - Real-world patterns
- [Docker Examples](./docker.md) - Container integration
- [Main README](../README.md) - Quick start guide
