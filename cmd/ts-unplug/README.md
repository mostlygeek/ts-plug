# ts-unplug

Bring Tailscale services to localhost - access remote services as if they were local.

## Build

```sh
cd ../..
make ts-unplug
```

## Usage

```sh
ts-unplug -dir ./state -port 8080 myserver.tailnet-name.ts.net
# Access at http://localhost:8080
```

## Documentation

See **[docs/ts-unplug.md](../../docs/ts-unplug.md)** for complete documentation.
