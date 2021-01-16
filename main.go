package main

import (
	"encoding/json"
	"log"
)

const (
	port = ":50051"
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

	for w := 0; w < 5; w++ {
		go func() {
			addr := "amqp://guest:guest@localhost:5672/"
			queue := New("task", addr)
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
	}

	go func() {
		if err := monitoring.Run(); err != nil {
			log.Fatalf("failed run monitoring: %v", err)
		}
	}()

	if err := RunGrpcServer(taskServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
