> [!WARNING]
> Lots of Work in Progress stuff here!

# Dude, Where's my Sidecar?!

- Turn a web server written in anything into a `tsnet.Server{}` with identity information.
- Also `docker` contains various no-side car docker containers I'm testing

## usage

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
