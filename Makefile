# Generete Code by Protocol Buffer
generate:
	sh generate.sh

# test go
test:
	docker-compose exec post_api go test -v ./grpc; \
	docker-compose exec post_api go test -v ./usecase/interactor

# test go
aws_test:
	docker-compose exec post_api go test -v ./aws;

# golint
lint:
	docker-compose exec post_api golint ./...
