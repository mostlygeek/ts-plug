> [!WARNING]
> Lots of Work in Progress stuff here!

# ts-plug

One line to turn a server written in anything into an application node on your tailnet!

```sh
$ ./ts-plug -hostname hello -- hello.js
```

## Table of Contents

- [Binaries Overview](#binaries-overview)
- [Quick Start](#quick-start)
- [Examples](#examples)
- [Docker Integration](#docker-integration)
- [Documentation](#documentation)

## Binaries Overview

This repository provides two complementary tools for Tailscale connectivity:

| Binary | Purpose | Direction | Use Case |
|--------|---------|-----------|----------|
| **ts-plug** | Expose localhost to Tailscale | Local → tailnet | Share your dev server with teammates |
| **ts-unplug** | Expose Tailscale to localhost | tailnet → Local | Access remote services locally |

### Architecture Diagram

```
ts-plug: Expose Local Service to tailnet
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│ Your Server  │  HTTP   │   ts-plug    │  HTTPS  │   tailnet    │
│ localhost:80 │ ──────> │ (w/ TLS)     │ ──────> │   Devices    │
└──────────────┘         └──────────────┘         └──────────────┘

ts-unplug: Expose tailnet Service to Local
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   tailnet    │  HTTPS  │  ts-unplug   │  HTTP   │ localhost:80 │
│   Service    │ ──────> │              │ ──────> │ (your apps)  │
└──────────────┘         └──────────────┘         └──────────────┘
```

### When to Use Each Tool

**Use `ts-plug` when:**
- Running a local dev server you want to share
- Testing webhooks that need a public URL
- Deploying services in containers without sidecars
- Sharing a local service with your team

**Use `ts-unplug` when:**
- Accessing a remote Tailscale service as if it's local
- Testing against a staging environment on your tailnet
- Developing against services that expect localhost
- Using tools that don't support HTTPS or custom domains

## Quick Start

### Building

Build both binaries:
```sh
make
```

Build platform-specific binaries:
```sh
make darwin  # macOS arm64
make linux   # Linux arm64 and amd64
```

Install to $GOPATH/bin:
```sh
make install
```

### ts-plug: Expose Local Service

Share a local web server on your tailnet:
```sh
./build/ts-plug -hostname myapp -- python -m http.server 8080
```

Make it publicly accessible:
```sh
./build/ts-plug -hostname myapp -public -- python -m http.server 8080
```

### ts-unplug: Access Remote Service

Bring a remote Tailscale service to localhost:
```sh
./build/ts-unplug -dir ./state -port 8080 myserver.tailnet.ts.net
# Now access at http://localhost:8080
```

Access a service on a specific port:
```sh
./build/ts-unplug -dir ./state -port 3000 database.tailnet.ts.net:5432
```

For detailed usage of each binary, see:
- [ts-plug Documentation](./cmd/ts-multi-plug/)
- [ts-unplug Documentation](./cmd/ts-unplug/README.md)

## Examples

The `cmd/examples/` directory contains example servers in multiple languages:

```sh
# Build examples
make examples

# Go server
./build/ts-plug -hn hello -- ./build/hello

# Node.js
./build/ts-plug -hn hello -- cmd/examples/hello/hello.js

# Python
./build/ts-plug -hn hello -- cmd/examples/hello/hello.py

# Ruby
./build/ts-plug -hn hello -- cmd/examples/hello/hello.rb

# Perl
./build/ts-plug -hn hello -- cmd/examples/hello/hello.pl

# Bash
./build/ts-plug -hn hello -- cmd/examples/hello/hello.sh
```

See [cmd/examples/](./cmd/examples/) for more details.

## ts-plug Features

ts-plug will automatically:

- start the upstream server
- join your tailnet
- generate a valid TLS cert and DNS name
- reverse proxy all traffic to http://127.0.0.1:8080

... and more:

- `-dns`: DNS reverse proxying (see docker/pi-hole)
- `-http`: http to http proxying
- `-public`: exposes your service to the public internet
- customize upstream ports with `-http-port`, `-https-port`, `-dns-port`,

## Make it Public

Use `-public` to share your server with everyone.

```sh
$ ./build/ts-plug -hn hello -public -- ./build/hello
```

It automatically provides a valid DNS name and TLS certificate, replacing the
need for other localhost tunneling solutions. Additionally, requests from your
tailnet include implicit identity information.

## Docker Integration

`ts-plug` removes the requirement for a Tailscale sidecar when running apps in
containers. Check out the examples in the `docker/` folder that show how to inject
`ts-plug` as the entrypoint:

- [Pi-hole](./docker/pi-hole/) - DNS-based ad blocker
- [Open WebUI](./docker/openwebui/) - ChatGPT-style web interface
- [Audiobookshelf](./docker/audiobookshelf/) - Audiobook and podcast server

It's still experimental but the initial results look good.

## Documentation

### Command-line Help

**ts-plug:**
```sh
./build/ts-plug -h
```

**ts-unplug:**
```sh
./build/ts-unplug -h
```

### Detailed Guides

- **[ts-plug Guide](./docs/ts-plug.md)** - Complete guide for exposing local services
- **[ts-unplug Guide](./docs/ts-unplug.md)** - Complete guide for accessing remote services
- **[Docker Examples](./docs/docker.md)** - Using ts-plug in containers
- **[Use Cases](./docs/use-cases.md)** - Real-world scenarios and patterns

### Quick Reference

**ts-plug flags:**
- `-hostname/-hn` - Hostname on tailnet (default: "tsmultiplug")
- `-dir` - Directory to store Tailscale state (default: ".data")
- `-http` - Enable HTTP listener (default: 80:8080)
- `-https` - Enable HTTPS listener (default: 443:8080)
- `-dns` - Enable DNS listener (default: 53:53)
- `-public` - Enable public HTTPS access (Tailscale Funnel)
- `-log` - Log level: debug, info, warn, error (default: "info")

**ts-unplug flags:**
- `-dir` - (required) Directory for tsnet server state
- `-hostname` - Hostname for the tsnet server (default: "tsunplug")
- `-port` - Local port to listen on (default: 80)
- `-debug-tsnet` - Enable tsnet.Server logging

## Contributing

Contributions welcome! This is a work-in-progress project exploring new patterns
for Tailscale integration.

## License

BSD-3-Clause - See [LICENSE](./LICENSE)
