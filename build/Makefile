.PHONY : build_protobuf run

build_protobuf:
	protoc -I=../proto --go_out=../internal/api --go_opt=paths=source_relative \
	--go-grpc_out=../internal/api --go-grpc_opt=paths=source_relative \
	../proto/auth.proto

run:
	go run ../internal/cmd/main.go