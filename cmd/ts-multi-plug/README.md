# ts-multi-plug

Reverse Proxies multiple services to the upstream child process

## Examples

### HTTPS reverse proxy

```sh
# these are all equilvent
$ ./ts-multi-plug -- server.js
$ ./ts-multi-plug -https -- server.js
$ ./ts-multi-plug -https -https-port 443:80 -- server.js
$ ./ts-multi-plug -name tsmultiplug -https -https-port 443:80 -- server.js
```

- automatically acquires a Let's Encrypt TLS certification
- available at https://tsmultiplug.ts-name.ts.net
- proxies all HTTP traffic to http://localhost:80

### HTTP reverse proxy

```sh
# these are all equilvent
$ ./ts-multi-plug -http -- server.js
$ ./ts-multi-plug -http -http-port 80:80 -- server.js
```

- available at http://tsmultiplug.ts-name.ts.net
- proxies all HTTP traffic to http://localhost:80

### DNS Reverse Proxy

```sh
# these are all equilvent
$ ./ts-multi-plug -dns -- resolver.js
$ ./ts-multi-plug -dns -dns-port 53:53 -- resolver.js

# alternative resolver port
$ ./ts-multi-plug -dns -dns-port 53:5353 -- resolver.js
```

- available at udp at tsmultiplug.ts-name.ts.net
- proxies all udp datagrams to udp://localhost:5353

### Multiple protocol proxying

```sh
$ ./ts-multi-plug -dns -https -- resolve-server.js
```
