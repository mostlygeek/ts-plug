# ts-plug

Plug any web server into tailscale with automatic identity information

## usage

```sh
$ make

# go web server
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
