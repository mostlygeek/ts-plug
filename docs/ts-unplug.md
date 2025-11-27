# ts-unplug Guide

**ts-unplug** is a reverse HTTP proxy that exposes a Tailscale service to localhost.

## Overview

ts-unplug allows you to:
- Access remote Tailscale services as if they were local
- Develop against remote APIs without changing code
- Test against staging environments seamlessly
- Use tools that only support localhost URLs

Think of it as the reverse of ts-plug: instead of exposing local to remote, it exposes remote to local.

## Installation

Build from source:
```sh
make ts-unplug
```

Install to $GOPATH/bin:
```sh
make install
```

## Basic Usage

```sh
ts-unplug -dir <state-dir> [options] <remote-addr>
```

### Simple Examples

Access a Tailscale service locally:
```sh
ts-unplug -dir ./state -port 8080 myserver.tailnet.ts.net
# Now access at http://localhost:8080
```

Connect to a specific port on a remote service:
```sh
ts-unplug -dir ./state -port 3000 database.tailnet.ts.net:5432
# PostgreSQL now available at localhost:3000
```

Access a staging API as localhost:
```sh
ts-unplug -dir ./state -port 8080 api-staging.tailnet.ts.net
# Your app can now use http://localhost:8080 for API calls
```

## Configuration Flags

### Required

- `<remote-addr>` - Remote Tailscale address (hostname or hostname:port)
  - If no port specified, defaults to port 80
  - Examples: `myserver`, `myserver.tailnet.ts.net`, `myserver:8080`

- `-dir` - Directory for tsnet server state (required)
  - Stores Tailscale authentication and connection state
  - Example: `-dir ./state`, `-dir /var/lib/tsunplug`

### Optional

- `-port` - Local port to listen on (default: 80)
  ```sh
  ts-unplug -dir ./state -port 8080 remote-api.tailnet.ts.net
  ```

- `-hostname` - Hostname for the tsnet server (default: "tsunplug")
  ```sh
  ts-unplug -dir ./state -hostname myproxy -port 3000 remote.tailnet.ts.net
  ```

- `-debug-tsnet` - Enable verbose tsnet.Server logging
  ```sh
  ts-unplug -dir ./state -port 8080 -debug-tsnet remote.tailnet.ts.net
  ```

## Use Cases

### Local Development Against Remote Services

Develop locally while using a remote database:
```sh
# Start the proxy
ts-unplug -dir ./state -port 5432 postgres.tailnet.ts.net:5432

# In another terminal, run your app pointing to localhost
DATABASE_URL=postgresql://localhost:5432/mydb npm run dev
```

### Testing Against Staging

Test your frontend against a staging API:
```sh
ts-unplug -dir ./state -port 8080 api-staging.tailnet.ts.net

# Update your .env.local
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local

npm run dev
```

### Tool Integration

Many tools only work with localhost URLs. Use ts-unplug to bridge the gap:

```sh
# Access a remote admin panel locally
ts-unplug -dir ./state -port 8080 admin.tailnet.ts.net

# Now use curl, Postman, etc. with localhost
curl http://localhost:8080/api/status
```

### Multiple Services

Run multiple instances to access different services:

```sh
# Terminal 1: Database
ts-unplug -dir ./state-db -port 5432 postgres.tailnet.ts.net:5432

# Terminal 2: Redis
ts-unplug -dir ./state-redis -port 6379 redis.tailnet.ts.net:6379

# Terminal 3: API
ts-unplug -dir ./state-api -port 8080 api.tailnet.ts.net
```

Note: Each instance needs its own `-dir` to avoid conflicts.

### Debugging Remote Services

Debug a remote service with local tools:
```sh
ts-unplug -dir ./state -port 8080 buggy-service.tailnet.ts.net

# Use your favorite debugging tools
curl -v http://localhost:8080/debug
http localhost:8080/health  # HTTPie
```

## How It Works

```
┌─────────────────────────────────────────────────────────┐
│  Your Local Machine                                     │
│                                                         │
│  ┌──────────────┐          ┌──────────────┐             │
│  │ Your App     │  HTTP    │  ts-unplug   │             │
│  │ localhost:80 │ ──────>  │              │             │
│  └──────────────┘          └──────┬───────┘             │
│                                   │                     │
└───────────────────────────────────┼─────────────────────┘
                                    │ Tailscale
                                    │ (encrypted)
┌───────────────────────────────────┼─────────────────────┐
│  Remote Tailscale Network         │                     │
│                                   │                     │
│                          ┌────────▼───────┐             │
│                          │ Remote Service │             │
│                          │ myserver:80    │             │
│                          └────────────────┘             │
└─────────────────────────────────────────────────────────┘
```

ts-unplug:
1. Connects to your tailnet
2. Establishes a connection to the remote service
3. Listens on localhost
4. Forwards all traffic through the encrypted Tailscale connection

## Advanced Usage

### Port Mapping

The remote and local ports don't need to match:
```sh
# Remote service on :8080, local access on :3000
ts-unplug -dir ./state -port 3000 api.tailnet.ts.net:8080
```

### Long-Running Proxy

Run as a background service:
```sh
# Using systemd (example)
cat > /etc/systemd/system/ts-unplug-api.service <<EOF
[Unit]
Description=ts-unplug proxy for API
After=network.target

[Service]
Type=simple
User=youruser
ExecStart=/usr/local/bin/ts-unplug -dir /var/lib/tsunplug -port 8080 api.tailnet.ts.net
Restart=always

[Install]
WantedBy=multi-user.target
EOF

systemctl enable --now ts-unplug-api
```

### Docker Container

Access a service from inside a Docker container:
```sh
# Start ts-unplug on host
ts-unplug -dir ./state -port 8080 api.tailnet.ts.net

# Run container with access to host network
docker run --network host myapp
# Container can now access http://localhost:8080
```

Or use host.docker.internal:
```sh
docker run -e API_URL=http://host.docker.internal:8080 myapp
```

## Security Considerations

### Authentication

ts-unplug inherits your Tailscale authentication. The remote service sees requests as coming from your Tailscale identity.

### State Directory

The `-dir` contains sensitive Tailscale credentials. Protect it appropriately:
```sh
chmod 700 ./state
```

### Localhost Binding

ts-unplug only listens on localhost (127.0.0.1), not on all network interfaces. This means only processes on your local machine can access it.

## Troubleshooting

### "dir is required" Error

You must specify a state directory:
```sh
# Wrong
ts-unplug myserver.tailnet.ts.net

# Right
ts-unplug -dir ./state myserver.tailnet.ts.net
```

### "remote-addr is required" Error

Provide the remote address as a positional argument:
```sh
ts-unplug -dir ./state -port 8080 myserver.tailnet.ts.net
```

### Connection Refused

If you can't connect to localhost:

1. Verify ts-unplug is running and shows "HTTP proxy listening"
2. Check you're using the correct local port
3. Verify the remote service is accessible from your Tailnet

### Can't Reach Remote Service

If ts-unplug starts but can't connect to the remote:

1. Verify the remote hostname is correct
2. Check the remote service is running
3. Ensure you have access to the remote service on your Tailnet
4. Try accessing the service directly: `tailscale ping myserver.tailnet.ts.net`

### Port Already in Use

If the local port is already taken:
```sh
# Check what's using the port
lsof -i :8080

# Use a different port
ts-unplug -dir ./state -port 8081 myserver.tailnet.ts.net
```

## Examples

### Access Remote PostgreSQL
```sh
ts-unplug -dir ./state -port 5432 postgres.tailnet.ts.net:5432

# Connect with psql
psql -h localhost -p 5432 -U myuser mydb
```

### Remote Redis
```sh
ts-unplug -dir ./state -port 6379 redis.tailnet.ts.net:6379

# Use redis-cli
redis-cli -h localhost -p 6379
```

### Remote HTTP API
```sh
ts-unplug -dir ./state -port 8080 api-staging.tailnet.ts.net

# Test with curl
curl http://localhost:8080/api/users
```

### Remote Admin Panel
```sh
ts-unplug -dir ./state -port 3000 admin.tailnet.ts.net:3000

# Open in browser
open http://localhost:3000
```

## Comparison with ts-plug

| Feature | ts-plug | ts-unplug |
|---------|---------|-----------|
| Direction | Local → tailnet | tailnet → Local |
| Use Case | Share local services | Access remote services |
| Starts Process | Yes | No |
| TLS | Automatic | Proxies existing |
| Public Access | Optional | No |

## See Also

- [ts-plug Guide](./ts-plug.md) - Expose local services to Tailnet
- [Use Cases](./use-cases.md) - Real-world patterns
- [Main README](../README.md) - Quick start guide
