FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

RUN go env -w CGO_ENABLED="1"

WORKDIR /build

RUN git clone https://github.com/anxin100/harvest_bot

RUN  cd /build/harvest_bot && go build -o /harvest_bot

FROM alpine

WORKDIR /root

COPY  --from=build /harvest_bot /usr/bin/harvest_bot

ENTRYPOINT [ "harvest_bot" ]