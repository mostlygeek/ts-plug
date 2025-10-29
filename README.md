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

## Dude, Where's my Sidecar?!

`ts-plug` removes the requirement for a Tailscale sidecar when running apps in
containers. Check out the examples in the `docker/` folder that show how to inject
`ts-plug` as the entrypoint. It's still experimental but the initial results look good.

## Help

Use the `-h` flag to get all the CLI flags available to tsplug

```sh
$ ./build/ts-plug -h
Usage of ./build/ts-plug:
  -debug-tsnet
        enable tsnet.Server logging
  -dir string
        directory to store tailscale state (default ".data")
  -dns
        Enable DNS listener (default 53:53)
  -dns-port value
        DNS port mapping (in:out or port) (default 53:53)
  -hn string
        hostname on tailnet (short) (default "tsmultiplug")
  -hostname string
        hostname on tailnet (default "tsmultiplug")
  -http
        Enable HTTP listener (default 80:8080)
  -http-port value
        HTTP port mapping (in:out or port) (default 80:8080)
  -https
        Enable HTTPS listener (default 443:8080)
  -https-port value
        HTTPS port mapping (in:out or port) (default 443:8080)
  -log string
        Log level (debug | info | warn | error) (default "info")
  -public
        Enable public https access
```
