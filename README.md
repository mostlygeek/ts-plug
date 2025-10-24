# ts-plug

Turn a web server written in anything into a `tsnet.Server{}` with identity information...

## usage

```sh
$ make

# go
$ ./build/ts-plug-web -- ./build/hello

# node
$ ./build/ts-plug-web -- cmd/examples/hello-node/hello.js

# python
$ ./build/ts-plug-web -- cmd/examples/hello-python/hello.py

# ruby
$ ./build/ts-plug-web -- cmd/examples/hello-ruby/hello.rb

# perl
$ ./build/ts-plug-web -- cmd/examples/hello-perl/hello.pl

# bash, can't forget the GOAT
$ ./build/ts-plug-web -- cmd/examples/hello-sh/hello.sh
```

Add `-funnel` to make it also accessible over the Internet (no identity though)

```sh
# go
$ ./build/ts-plug-web -funnel -- ./build/hello

# node
$ ./build/ts-plug-web -funnel -- cmd/examples/hello-node/hello.js

# ... etc.
```
