> [!WARNING]
> Lots of Work in Progress stuff here!

# ts-plug

One line to turn a server written in anything into an application node on your tailnet!

```sh
./ts-plug -hostname hello -- hello.js
```

## What's in this repo?

| Binary | Purpose | Use Case |
|--------|---------|----------|
| **ts-plug** | Expose localhost to your tailnet | Share your dev server, deploy without sidecars |
| **ts-unplug** | Bring tailnet services to localhost | Access remote databases/APIs as if they were local |

## Quick Start

**Build:**
```sh
make                    # Build both binaries
make install            # Install to $GOPATH/bin
```

**ts-plug** - Share a local service:
```sh
./build/ts-plug -hostname myapp -- python -m http.server 8080
# Access at https://myapp.tailnet-name.ts.net
```

**ts-unplug** - Access a remote service:
```sh
./build/ts-unplug -dir ./state -port 8080 api.tailnet-name.ts.net
# Access at http://localhost:8080
```

## Key Features

**ts-plug** automatically:
- Starts your upstream server
- Joins your tailnet with TLS and DNS
- Reverse proxies to localhost:8080
- Optional public access with `-public`
- Supports HTTP, HTTPS, and DNS protocols

**ts-unplug** provides:
- Reverse proxy from tailnet to localhost
- Access to services requiring localhost URLs
- Simple port mapping

## Examples

Run servers in any language:
```sh
make examples

# Try different languages with ts-plug
./build/ts-plug -hn hello -- ./build/hello        # Go
./build/ts-plug -hn hello -- cmd/examples/hello/hello.js   # Node
./build/ts-plug -hn hello -- cmd/examples/hello/hello.py   # Python
```

See [cmd/examples/](./cmd/examples/) for more.

## Docker Integration

Use ts-plug as an entrypoint to eliminate sidecar containers:

```dockerfile
COPY ts-plug /usr/local/bin/
ENTRYPOINT ["ts-plug", "-hostname", "myapp", "--"]
CMD ["npm", "start"]
```

See [docker/](./docker/) for Pi-hole, Open WebUI, and Audiobookshelf examples.

## Documentation

- **[Complete Documentation](./docs/)** - Guides, use cases, and detailed examples
- **[ts-plug Guide](./docs/ts-plug.md)** - Full ts-plug documentation
- **[ts-unplug Guide](./docs/ts-unplug.md)** - Full ts-unplug documentation
- **[Use Cases](./docs/use-cases.md)** - Real-world scenarios
- **[Docker Guide](./docs/docker.md)** - Container integration

**Quick help:**
```sh
./build/ts-plug -h
./build/ts-unplug -h
```

## License

BSD-3-Clause - See [LICENSE](./LICENSE)
