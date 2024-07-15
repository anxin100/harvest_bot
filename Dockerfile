FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

RUN go env -w CGO_ENABLED="1"

WORKDIR /build

RUN git clone https://code.wuban.net.cn/odysseus/aon-app-server

RUN  cd /build/aon-app-server && go build -o /aon-app-server

FROM alpine

WORKDIR /root

COPY  --from=build /aon-app-server /usr/bin/aon-app-server

ENTRYPOINT [ "aon-app-server" ]