# ts-unplug

A reverse HTTP proxy that exposes a Tailscale service to localhost.

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
