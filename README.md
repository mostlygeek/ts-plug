> [!WARNING]
> Lots of Work in Progress stuff here!

# ts-multi-plug

Turn a server written in anything into a `tsnet.Server{}`.

## Dude, Where's my Sidecar?!

The Docker sidecar pattern using `tailscale/tailscale` is the most common pattern to host applications on a tailnet. This repo explores a new no-sidecar approach. What we do differently is:

1. Introduce `ts-multi-plug`, a `tsnet.Server{}` reverse proxy that also spawns a child server process.
   - `ts-multi-plug -hostname hello -https -- hello.js` - spawns a nodejs web server that is available at https://hello.my-ts.ts.net
2. The `ENTRYPOINT` is `ts-multi-plug`. It becomes the init process!

## usage (OUTDATED)

```sh
$ make

# go
$ ./build/ts-plug -- ./build/hello

# node
$ ./build/ts-plug -- cmd/examples/hello-node/hello.js

# python
$ ./build/ts-plug -- cmd/examples/hello-python/hello.py

# ruby
$ ./build/ts-plug -- cmd/examples/hello-ruby/hello.rb

# perl
$ ./build/ts-plug -- cmd/examples/hello-perl/hello.pl

# bash ... because we can!
$ ./build/ts-plug -- cmd/examples/hello-sh/hello.sh
```

Add `-funnel` to make it also accessible over the Internet (no identity though)

```sh
# go
$ ./build/ts-plug -funnel -- ./build/hello

# node
$ ./build/ts-plug -funnel -- cmd/examples/hello-node/hello.js

# ... etc.
```
