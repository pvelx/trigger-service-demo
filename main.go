package main

import (
	"encoding/json"
	"github.com/pvelx/triggerHook/domain"
	"github.com/streadway/amqp"
	"log"
)

const (
	port = ":50051"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Dial: %s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error open channel:%v", err)
	}

	transport := func(task domain.Task) {
		taskJson, err := json.Marshal(task)
		if err != nil {
			log.Fatal(err)
		}

		if err = channel.Publish(
			"test-exchange", // publish to an exchange
			"",              // routing to 0 or more queues
			true,            // mandatory
			false,           // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "application/json",
				ContentEncoding: "",
				Body:            taskJson,
				DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
				Priority:        0,              // 0-9
			},
		); err != nil {
			log.Fatalf("Exchange Publish: %s", err)
		}
	}

	monitoring := NewMonitoring()
	tasksDeferredService := BuildTriggerHook(monitoring, transport)
	taskServer := NewTaskServer(tasksDeferredService)

	go func() {
		if err := tasksDeferredService.Run(); err != nil {
			log.Fatalf("failed run trigger hook: %v", err)
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
