FROM golang:1.23.7-alpine as builder

ENV CGO_ENABLED=0

WORKDIR /src

COPY . /src

# build the agent (with two different versions, so we can test auto-update mechanism)
RUN --mount=type=cache,target=/go \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -ldflags "-s -w" \
    -o /usr/bin/qbee-http-broker

FROM alpine:3.20.3

COPY --from=builder /usr/bin/qbee-http-broker /usr/bin/qbee-http-broker

EXPOSE 8081/tcp

ENTRYPOINT ["/usr/bin/qbee-http-broker" ]

