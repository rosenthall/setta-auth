.PHONY : build_protobuf run

build_protobuf:
	protoc -I=../proto --go_out=../internal/api --go_opt=paths=source_relative \
	--go-grpc_out=../internal/api --go-grpc_opt=paths=source_relative \
	../proto/auth.proto

build:
	CGO_ENABLED=0 go build  -o target/auth ../internal/cmd/main.go

generate-keys:
	@mkdir -p ../keys
	@openssl genpkey -algorithm RSA -out ../keys/private_key.pem -pkeyopt rsa_keygen_bits:2048
	@openssl rsa -in ../keys/private_key.pem -pubout -out ../keys/public_key.pem
	@echo "Keys generated in keys/ directory"

# Only for local development
setup-redis:
	sudo -S redis-server ../../infra/redis/auth_service/redis.conf

run:
	go run ../internal/cmd/main.go
