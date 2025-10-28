FROM pihole/pihole:2025.10.3
ARG TARGETARCH
COPY ../../build/ts-multi-plug-linux-${TARGETARCH} /usr/local/bin/ts-multi-plug
RUN chmod +x /usr/local/bin/ts-multi-plug
ENTRYPOINT ["/usr/local/bin/ts-multi-plug" "-dir", "/var/run/tsnet", "-hostname", "tsplug-pihole", "-https", "-https-port", "443:80", "-dns", "--"]
CMD ["/usr/bin/start.sh"]