# ts-plug

> **Note:** This is the source directory for the `ts-plug` binary. For comprehensive documentation, see the [main README](../../README.md) and [detailed guide](../../docs/ts-plug.md).

## Quick Links

- **[Main README](../../README.md)** - Overview of both ts-plug and ts-unplug
- **[Complete ts-plug Guide](../../docs/ts-plug.md)** - Detailed documentation
- **[ts-unplug](../ts-unplug/)** - Companion tool (reverse direction)
- **[Use Cases](../../docs/use-cases.md)** - Real-world patterns
- **[Docker Guide](../../docs/docker.md)** - Container integration

## Quick Start

Build:
```sh
cd ../..  # Return to repository root
make ts-plug
```

Basic usage:
```sh
./build/ts-plug -hostname myapp -- python -m http.server 8080
```

View help:
```sh
./build/ts-plug -h
```

## What is ts-plug?

ts-plug exposes local services to your Tailscale network:
- Wraps your application and handles Tailscale connectivity
- Provides automatic HTTPS with valid TLS certificates
- Optional public access via Tailscale Funnel
- Supports HTTP, HTTPS, and DNS proxying

## Common Examples

### Development Server
```sh
ts-plug -hostname myapp -- npm run dev
```

### Public Demo
```sh
ts-plug -public -hostname demo -- python app.py
```

### Multiple Protocols
```sh
ts-plug -http -https -hostname web -- ./server
```

## Source Code

The main implementation is in [`ts-multi-plug.go`](./ts-multi-plug.go).

## See Also

- **ts-unplug** - Reverse proxy to access remote Tailscale services locally
- **Examples** - Sample servers in [multiple languages](../examples/)
- **Docker** - Container deployment [examples](../../docker/)
