FROM golang:1.21.3-bookworm AS builder

COPY . /src
WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags \
     '-s -w -extldflags "-static"'  -o /esmqtt main.go

FROM alpine:3.17

RUN apk add --no-cache curl

COPY --from=builder /esmqtt /usr/local/bin/esmqtt

RUN chmod +x /usr/local/bin/esmqtt

ENTRYPOINT ["/usr/local/bin/esmqtt"]