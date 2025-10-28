FROM ghcr.io/advplyr/audiobookshelf:2.30.0
ARG TARGETARCH
COPY ../../build/ts-multi-plug-linux-${TARGETARCH} /app/ts-multi-plug
RUN chmod +x /app/ts-multi-plug
ENTRYPOINT ["/app/ts-multi-plug", "-dir", "/var/run/tsnet", "-hostname", "tsplug-audiobookshelf", "-https", "-https-port", "443:80", "--"]
CMD ["node", "index.js"]