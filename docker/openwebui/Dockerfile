FROM ghcr.io/open-webui/open-webui:slim
ARG TARGETARCH
COPY ../../build/ts-multi-plug-linux-${TARGETARCH} /usr/local/bin/ts-multi-plug
RUN chmod +x /usr/local/bin/ts-multi-plug
ENTRYPOINT ["/usr/local/bin/ts-multi-plug", "-https", "--"]
CMD ["/app/backend/start.sh"]