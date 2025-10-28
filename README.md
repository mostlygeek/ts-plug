> [!WARNING]
> Lots of Work in Progress stuff here!

# ts-plug

One line to turn a server written in anything into an application node on your Tailnet!

```sh
$ ./ts-plug -hostname hello -- hello.js
```

## More examples

```sh
# build binaries
$ make

# go
$ ./build/ts-plug -hn hello -- ./build/hello

# node
$ ./build/ts-plug -hn hello -- cmd/examples/hello-node/hello.js

# python
$ ./build/ts-plug -hn hello -- cmd/examples/hello-python/hello.py

# ruby
$ ./build/ts-plug -hn hello -- cmd/examples/hello-ruby/hello.rb

# perl
$ ./build/ts-plug -hn hello -- cmd/examples/hello-perl/hello.pl

# bash ... but of course!
$ ./build/ts-plug -hn hello -- cmd/examples/hello-sh/hello.sh

```

## Funnel Support

Make your application node available to everyone with [Funnel](https://tailscale.com/kb/1223/funnel)!

Add `-funnel` to make it also accessible over the Internet (no identity though)

```sh

# !!! Tip !!
# Try accessing this with Tailscale connected and disconnected. Your
# identity is automatically available in to the hello server

$ ./build/ts-plug -hn hello -funnel -- ./build/hello
```

## Dude, Where's my Sidecar?!

Using `ts-plug` it is possible to remove the need for a tailscale sidecar
when running containerized applications. In the `docker/` folder are examples
injecting in `ts-plug` and having it be the entrypoint. It is still very
experimental but initial experiments are positive.
