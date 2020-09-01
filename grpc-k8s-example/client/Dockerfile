FROM golang:1.15-alpine AS build
WORKDIR /src
ARG GOARCH
ARG GOARM
COPY . .
RUN CGO_ENABLED=0 GOARCH=$GOARCH GOARM=$GOARM go build -o /build/client
FROM alpine:3.12.0
COPY --from=build /build/client /bin/client
ENTRYPOINT ["/bin/client"]
