FROM pihole/pihole:latest
ARG TARGETARCH
COPY ../../build/ts-multi-plug-linux-${TARGETARCH} /usr/local/bin/ts-multi-plug
RUN chmod +x /usr/local/bin/ts-multi-plug
ENTRYPOINT ["/usr/local/bin/ts-multi-plug"]
CMD ["-dir", "/var/run/tsnet", "-hostname", "tsplug-dns2", "-https", "-https-port", "443:80", "-dns", "--", "/usr/bin/start.sh"]