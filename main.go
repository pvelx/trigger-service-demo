package main

import (
	"github.com/pvelx/triggerHook"
	"github.com/pvelx/triggerHookExample/proto"
	"github.com/pvelx/triggerHookExample/task_server"
	"google.golang.org/grpc"
	"log"
	"net"
)

var tasksDeferredService = triggerHook.Default()

const (
	port = ":50051"
)

func main() {
	tasksDeferredService.SetTransport(NewTransportAmqp())
	go tasksDeferredService.Run()

	taskServer := task_server.New(tasksDeferredService)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterTaskServer(s, taskServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
