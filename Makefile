gen_proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/task.proto

build:
	go build -o main github.com/pvelx/trigger-service-demo