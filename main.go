package main

import (
	"encoding/json"
	"log"
)

const (
	port     = ":50051"
	rabbitMq = "amqp://guest:guest@localhost:5672/"
	queue    = "task"
)

func main() {
	monitoring := NewMonitoring()
	tasksDeferredService := BuildTriggerHook(monitoring)
	taskServer := NewTaskServer(tasksDeferredService)

	go func() {
		if err := tasksDeferredService.Run(); err != nil {
			log.Fatalf("failed run trigger hook: %v", err)
		}
	}()

	queue := New(queue, rabbitMq)
	go func() {
		for {
			result := tasksDeferredService.Consume()
			taskJson, err := json.Marshal(result.Task())
			if err != nil {
				log.Fatal(err)
			}

			if err := queue.Push(taskJson); err != nil {
				log.Println("the task was not sent")
				result.Rollback()
			}

			result.Confirm()
		}
	}()

	go func() {
		if err := monitoring.Run(); err != nil {
			log.Fatalf("failed run monitoring: %v", err)
		}
	}()

	if err := RunGrpcServer(taskServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
