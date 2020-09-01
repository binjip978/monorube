proto: service/service.proto
	protoc -I=service --go_out=plugins=grpc:service service/service.proto

build-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

docker-amd64:
	docker build --build-arg GOARCH=amd64 -t binjip978/grpc-example-server-amd64:0.0.1 .

docker-arm7:
	docker build --build-arg GOARCH=arm --build-arg GOARM=7 -t binjip978/grpc-example-server-arm7:0.0.1 .

docker-push-amd64: docker-amd64
	docker push binjip978/grpc-example-server-amd64:0.0.1

docker-push-arm7: docker-arm7
	docker push binjip978/grpc-example-server-arm7:0.0.1
