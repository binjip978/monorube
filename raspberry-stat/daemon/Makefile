.PHONY: build test

build: test
	@GOARCH=arm GOOS=linux go build

test:
	@go test -race ./...

docker: build
	docker build . -t binjip978/rasperrypi-temp:v1.0.0

