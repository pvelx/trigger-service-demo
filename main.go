package main

import (
	"encoding/json"
	"github.com/pvelx/triggerHook/connection"
	"log"
	"os"
)

var (
	host     = os.Getenv("SERVER_HOST")
	rabbitMq = os.Getenv("MESSENGER_TRANSPORT_DSN")
	queue    = os.Getenv("QUEUE_NAME")
	exchange = os.Getenv("EXCHANGE_NAME")

	mysqlUser     = os.Getenv("DATABASE_USER")
	mysqlPassword = os.Getenv("DATABASE_PASSWORD")
	mysqlHost     = os.Getenv("DATABASE_HOST")
	mysqlDbName   = os.Getenv("DATABASE_NAME")
)

func main() {
	monitoring := NewMonitoring()
	tasksDeferredService := BuildTriggerHook(monitoring, connection.Options{
		User:     mysqlUser,
		Password: mysqlPassword,
		Host:     mysqlHost,
		DbName:   mysqlDbName,
	})
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
