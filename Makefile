run:
	go run cmd/main.go
proto:
	protoc --go_out=./ --go-grpc_out=./  ./proto/proto.proto