package main

import (
	"fmt"
	"github.com/pvelx/triggerHook"
	"github.com/pvelx/triggerHook/connection"
	"github.com/pvelx/triggerHook/contracts"
	"github.com/pvelx/triggerHook/domain"
	"github.com/pvelx/triggerHook/error_service"
	"github.com/pvelx/triggerHook/monitoring_service"
	"github.com/pvelx/triggerHook/sender_service"
	"path"
	"time"
)

func BuildTriggerHook(monitoring *Monitoring, transport func(domain.Task)) contracts.TasksDeferredInterface {

	eventHandlers := make(map[contracts.Level]func(event contracts.EventError))
	baseFormat := "%s MESSAGE:%s METHOD:%s FILE:%s:%d EXTRA:%v\n"
	for level, format := range map[contracts.Level]string{
		contracts.LevelDebug: "DEBUG:" + baseFormat,
		contracts.LevelError: "ERROR:" + baseFormat,
		contracts.LevelFatal: "FATAL:" + baseFormat,
	} {
		format := format
		level := level
		eventHandlers[level] = func(event contracts.EventError) {
			_, shortMethod := path.Split(event.Method)
			_, shortFile := path.Split(event.File)
			fmt.Printf(
				format,
				event.Time.Format("2006-01-02 15:04:05.000"),
				event.EventMessage,
				shortMethod,
				shortFile,
				event.Line,
				event.Extra,
			)
		}
	}

	subscriptions := make(map[contracts.Topic]func(event contracts.MeasurementEvent))
	for _, topic := range []contracts.Topic{
		contracts.SpeedOfPreloading,
		contracts.SpeedOfCreating,
		contracts.CountOfAllTasks,
		contracts.Preloaded,
		contracts.SpeedOfConfirmation,
		contracts.CountOfWaitingForSending,
		contracts.WaitingForConfirmation,
		contracts.SpeedOfSending,
		contracts.SpeedOfDeleting,
	} {
		topic := topic
		subscriptions[topic] = func(event contracts.MeasurementEvent) {
			monitoring.AddMeasurement(string(topic), event)
		}
	}

	tasksDeferredService := triggerHook.Build(triggerHook.Config{
		Connection: connection.Options{
			User:         "root",
			Password:     "secret",
			Host:         "127.0.0.1:3306",
			DbName:       "test_db",
			MaxOpenConns: 30,
			MaxIdleConns: 30,
		},
		ErrorServiceOptions: error_service.Options{
			Debug:         false,
			EventHandlers: eventHandlers,
		},
		MonitoringServiceOptions: monitoring_service.Options{
			PeriodMeasure: 5 * time.Second,
			Subscriptions: subscriptions,
		},
		SenderServiceOptions: sender_service.Options{
			Transport: transport,
		},
	})

	return tasksDeferredService
}
