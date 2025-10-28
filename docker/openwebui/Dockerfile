FROM ghcr.io/open-webui/open-webui:slim
ARG TARGETARCH
COPY ../../build/ts-plug-linux-${TARGETARCH} /usr/local/bin/ts-plug
RUN chmod +x /usr/local/bin/ts-plug
ENTRYPOINT ["/usr/local/bin/ts-plug"]
CMD ["--", "/app/backend/start.sh"]