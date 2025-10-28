FROM ghcr.io/advplyr/audiobookshelf:2.30.0
ARG TARGETARCH
COPY ../../build/ts-plug-linux-${TARGETARCH} /app/ts-plug
RUN chmod +x /app/ts-plug
ENTRYPOINT ["/app/ts-plug", "-dir", "/var/run/tsnet", "-hostname", "tsplug-audiobookshelf", "-https", "-https-port", "443:80", "--"]
CMD ["node", "index.js"]