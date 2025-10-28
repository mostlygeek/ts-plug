# Fake DNS Resolver for Testing

A simple DNS server that responds to all queries with fixed values. Useful for testing `ts-plug -dns`.

## Supported Record Types

The resolver supports the following DNS record types and always returns the same values:

- **A**: `192.0.2.1` (TEST-NET-1 address)
- **AAAA**: `2001:db8::1` (documentation IPv6 address)
- **CNAME**: `www.<domain>` resolves to `<domain>`
- **TXT**: `v=test dns resolver`
- **MX**: `10 mx.<domain>` (priority 10)

Where `<domain>` defaults to `tailscale.com` but can be customized with the `-domain` flag.

## Usage

Build and run the resolver:

```bash
go run cmd/examples/resolver/resolver.go
```

Or specify a custom port and domain:

```bash
go run cmd/examples/resolver/resolver.go -port 5353 -domain example.com
```

Available flags:

- `-port`: Port to listen on (default: `53`)
- `-domain`: Domain to use for CNAME and MX responses (default: `tailscale.com`)

## Testing with ts-plug

1. Start the fake DNS resolver on a local port:

   ```bash
   go run cmd/ts-multi-plug/ts-multi-plug.go -dns -- go run cmd/examples/resolver/resolver.go
   ```

2. From another machine on your tailnet, test DNS queries:
   ```bash
   dig @your-hostname.tailnet.ts.net tailscale.com A
   dig @your-hostname.tailnet.ts.net tailscale.com AAAA
   dig @your-hostname.tailnet.ts.net www.tailscale.com CNAME
   dig @your-hostname.tailnet.ts.net tailscale.com TXT
   dig @your-hostname.tailnet.ts.net tailscale.com MX
   ```

## Example Output

```
$ go run cmd/examples/resolver/resolver.go -port 5353 -domain tailscale.com
15:30:45.123456 Fake DNS resolver listening on 127.0.0.1:5353
15:30:45.123456 Resolving all queries to fixed values:
15:30:45.123456   A:     192.0.2.1
15:30:45.123456   AAAA:  2001:db8::1
15:30:45.123456   CNAME: www.tailscale.com -> tailscale.com
15:30:45.123456   TXT:   v=test dns resolver
15:30:45.123456   MX:    10 mx.tailscale.com
15:30:50.234567 Query from 127.0.0.1:54321: tailscale.com (type 1)
15:30:51.345678 Query from 127.0.0.1:54322: www.tailscale.com (type 5)
```
