package main

import (
	"context"
	"github.com/pvelx/trigger-service-demo/proto"
	"github.com/pvelx/triggerhook/contracts"
	"github.com/pvelx/triggerhook/domain"
)

func NewTaskServer(tasksDeferredService contracts.TriggerHookInterface) proto.TaskServer {
	return &taskServer{
		tasksDeferredService: tasksDeferredService,
	}
}

type taskServer struct {
	proto.UnimplementedTaskServer
	tasksDeferredService contracts.TriggerHookInterface
}

func (s *taskServer) Create(ctx context.Context, req *proto.Request) (*proto.Response, error) {

	task := domain.Task{
		ExecTime: req.ExecTime,
	}
	if err := s.tasksDeferredService.Create(&task); err != nil {
		return &proto.Response{Status: "fail"}, nil
	}

	return &proto.Response{Status: "ok", Id: task.Id}, nil
}

func (s *taskServer) Delete(ctx context.Context, req *proto.Request) (*proto.Response, error) {

	if err := s.tasksDeferredService.Delete(req.Id); err != nil {
		return &proto.Response{Status: "fail"}, nil
	}

	return &proto.Response{Status: "ok"}, nil
}
