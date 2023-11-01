FROM alpine:3.17

WORKDIR /var/esmqtt

COPY release/esmqtt /usr/local/bin/esmqtt

RUN chmod +x /usr/local/bin/esmqtt

ENTRYPOINT ["/usr/local/bin/esmqtt"]
