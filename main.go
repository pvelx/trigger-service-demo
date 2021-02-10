package main

import (
	"encoding/json"
	"github.com/pvelx/triggerhook/connection"
	"log"
	"os"
)

var (
	host     = os.Getenv("SERVER_HOST")
	rabbitMq = os.Getenv("MESSENGER_TRANSPORT_DSN")
	exchange = os.Getenv("EXCHANGE_NAME")

	mysqlUser     = os.Getenv("DATABASE_USER")
	mysqlPassword = os.Getenv("DATABASE_PASSWORD")
	mysqlHost     = os.Getenv("DATABASE_HOST")
	mysqlDbName   = os.Getenv("DATABASE_NAME")

	influxDbDns      = os.Getenv("INFLUX_DB_DNS")
	influxDbUsername = os.Getenv("INFLUX_DB_USERNAME")
	influxDbPassword = os.Getenv("INFLUX_DB_PASSWORD")
	influxDbName     = os.Getenv("INFLUX_DB_NAME")
)

func main() {
	monitoring := NewMonitoring(influxDbDns, influxDbUsername, influxDbPassword, influxDbName)
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

	queue := New(exchange, rabbitMq)
	go func() {
		for {
			result := tasksDeferredService.Consume()
			message := struct {
				TaskId string `json:"taskId"`
			}{
				result.Task().Id,
			}

			taskJson, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
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
