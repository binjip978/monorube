build-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

build-arm7:
	GOOS=linux GOARCH=arm GOARM=7 go build

docker-amd64: build-amd64
	docker build . -t binjip978/grpc-example-client-amd64:0.0.1

docker-arm7: build-arm7
	docker build . -t binjip978/grpc-example-client-arm7:0.0.1

docker-push-amd64: docker-amd64
	docker push binjip978/grpc-example-client-amd64:0.0.1

docker-push-arm7: docker-arm7
	docker push binjip978/grpc-example-client-arm7:0.0.1
