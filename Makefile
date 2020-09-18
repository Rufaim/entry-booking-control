all: build

build: build_proto build_server build_client

build_client:
	GOOS=linux go build -o entry_booking_client ./cmd/client

build_server:
	GOOS=linux go build -o entry_booking_server ./cmd/server

build_proto: 
	protoc --go_out=. --go-grpc_out=. ./cmd/message/message.proto