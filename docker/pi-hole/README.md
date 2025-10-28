## Quick notes

```sh
# build container
$ docker buildx build --platform linux/amd64,linux/arm64 -t mostlygeek/ts-plug:pihold-latest -f docker/pi-hole/Dockerfile --load .

# run container
$ docker run -it --rm --name "tsplug-dns" -v tsplug-dns:/var/run/tsnet mostlygeek/ts-plug:pihold-latest
```
