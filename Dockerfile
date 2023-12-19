FROM golang:1.20.0 as builder
MAINTAINER lanyulei <fdevops@163.com>

WORKDIR /opt/comet
COPY . .

RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy && make build

FROM alpine:3.17.1 as worker

WORKDIR /opt/comet
COPY --from=builder /opt/comet/comet /opt/comet/comet
COPY --from=builder /opt/comet/config/settings.yml /opt/comet/config

EXPOSE 8000

ENTRYPOINT ["/opt/comet/comet", "server", "-c", "/opt/comet/config/settings.yml"]
