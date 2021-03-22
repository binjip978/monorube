FROM golang:1.16.2 AS builder
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GO111MODULE=off go build -o hello

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/hello .
ENTRYPOINT ["/app/hello"]
