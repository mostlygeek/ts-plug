# ts-plug

Expose local services to your Tailscale network with automatic TLS and DNS.

## Build

```sh
cd ../..
make ts-plug
```

## Usage

```sh
ts-plug -hostname myapp -- python -m http.server 8080
# Access at https://myapp.tailnet-name.ts.net
```

## Documentation

See **[docs/ts-plug.md](../../docs/ts-plug.md)** for complete documentation.
