.PHONY: build
build:
	@go build

.PHONY: build-arm
build-arm:
	@GOARCH=arm GOOS=linux go build

.PHONY: clean
clean:
	@rm serial-temp-reader || true
