build: test
	go build

test:
	go test -cover ./...

docker: build
	docker build . -t binjip978/prime-service:latest

docker-run:
	docker run -d -p 8080:8080 binjip978/prime-service:latest

docker-push:
	docker push binjip978/prime-service:latest 

clean:
	rm -f ./prime-service

.PHONY: test build docker docker-run docker-push clean 