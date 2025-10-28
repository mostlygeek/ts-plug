## Building multi-platform container

```sh
# run in root of repo so the build/ts-plug-linux-{platform} can be copied in
# this is prototype so
$ buildx build --platform linux/amd64,linux/arm64 -t mostlygeek/ts-plug:owebui-slim -f docker/openwebui/Dockerfile --push .
```

## Running it

```sh
$docker run -it --rm  --add-host=host.docker.internal:host-gateway -v open-webui:/app/backend/data --name open-webui openwebui-plug:slim
```
