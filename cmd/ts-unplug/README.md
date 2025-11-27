# ts-unplug

> **Note:** This is the source directory for the `ts-unplug` binary. For comprehensive documentation, see the [main README](../../README.md) and [detailed guide](../../docs/ts-unplug.md).

## Quick Links

- **[Main README](../../README.md)** - Overview of both ts-plug and ts-unplug
- **[Complete ts-unplug Guide](../../docs/ts-unplug.md)** - Detailed documentation
- **[ts-plug](../ts-multi-plug/)** - Companion tool (reverse direction)
- **[Use Cases](../../docs/use-cases.md)** - Real-world patterns

## What is ts-unplug?

A reverse HTTP proxy that exposes a Tailscale service to localhost:
- Brings remote services to your local machine
- Perfect for development against remote APIs/databases
- Access services that require localhost URLs
- No need to modify application configuration

## Build

```sh
cd ../..  # Return to repository root
make ts-unplug
```

## Usage

```bash
ts-unplug -dir <state-dir> [options] <remote-addr>
```

## Arguments

- `<remote-addr>` - Remote Tailscale address to proxy (hostname or hostname:port). If no port is specified, defaults to port 80.

## Options

- `-dir <path>` - (required) Directory for tsnet server state
- `-hostname <name>` - Hostname for the tsnet server (default: "tsunplug")
- `-port <number>` - Local port to listen on (default: 80)
- `-debug-tsnet` - Enable tsnet.Server logging

## Example

Expose a Tailscale service running on `myserver` to `localhost:8080`:

```bash
ts-unplug -dir ./state -port 8080 mymachine.ts-name.ts.net
```

Access a service on a specific HTTP port:

```bash
ts-unplug -dir ./state -port 8080 myserver:3000
```

Once running, connect to the service at `http://localhost:<port>`.

## More Examples

### Remote Database
```bash
ts-unplug -dir ./state -port 5432 postgres.tailnet.ts.net:5432
psql -h localhost -p 5432 -U user dbname
```

### Remote API
```bash
ts-unplug -dir ./state -port 8080 api-staging.tailnet.ts.net
curl http://localhost:8080/api/endpoint
```

### Remote Redis
```bash
ts-unplug -dir ./state -port 6379 redis.tailnet.ts.net:6379
redis-cli -h localhost -p 6379
```

## Source Code

The main implementation is in [`ts-unplug.go`](./ts-unplug.go).

## See Also

- **[Complete Guide](../../docs/ts-unplug.md)** - Full documentation with examples
- **[ts-plug](../ts-multi-plug/)** - Expose local services to Tailscale
- **[Use Cases](../../docs/use-cases.md)** - Real-world patterns
- **[Main README](../../README.md)** - Project overview
