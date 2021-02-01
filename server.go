package main

import (
	"github.com/pvelx/triggerServiceDemo/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func RunGrpcServer(taskServer proto.TaskServer) error {

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	proto.RegisterTaskServer(server, taskServer)

	return server.Serve(lis)
}
