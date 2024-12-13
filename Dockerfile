# FROM golang:1.17-alpine

# RUN apk add ca-certificates

# WORKDIR /go/src/github.com/pierre-emmanuelJ/iptv-proxy
# COPY . .
# RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iptv-proxy .

# FROM alpine:3
# COPY --from=0  /go/src/github.com/pierre-emmanuelJ/iptv-proxy/iptv-proxy /
# ENTRYPOINT ["/iptv-proxy"]

FROM golang:1.23-alpine3.20 AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /go/src/github.com/pierre-emmanuelJ/iptv-proxy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iptv-proxy .

FROM alpine:3.20
COPY --from=builder  /go/src/github.com/pierre-emmanuelJ/iptv-proxy/iptv-proxy /
ENTRYPOINT ["/iptv-proxy"]
