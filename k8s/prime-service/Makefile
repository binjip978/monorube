VERSION :=

build: test
	go build

test:
	go test -cover ./...

docker: build version
	docker build . -t binjip978/prime-service:$(VERSION)

docker-run: version
	docker run -d -p 8080:8080 binjip978/prime-service:$(VERSION)

docker-push: version
	docker push binjip978/prime-service:$(VERSION)

clean:
	rm -f ./prime-service

version:
	$(eval VERSION = $(shell grep 'const Version' version.go | awk '{print $$4}'))
	@echo $(VERSION)

.PHONY: test build docker docker-run docker-push clean version
