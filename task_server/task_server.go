package task_server

import (
	"context"
	"github.com/pvelx/triggerHook/contracts"
	"github.com/pvelx/triggerHookExample/proto"
)

func New(tasksDeferredService contracts.SchedulerInterface) proto.TaskServer {
	return &taskServer{
		tasksDeferredService: tasksDeferredService,
	}
}

type taskServer struct {
	tasksDeferredService contracts.SchedulerInterface
	proto.UnimplementedTaskServer
}

func (s *taskServer) Create(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	task, e := s.tasksDeferredService.Create(req.ExecTime)
	if e != nil {
		return &proto.Response{Status: "fail"}, nil
	}
	//log.Printf("Created: %v", task.Id)

	return &proto.Response{Status: "ok", Id: task.Id}, nil
}
