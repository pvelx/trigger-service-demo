package main

import (
	"fmt"
	"github.com/pvelx/triggerhook"
	"github.com/pvelx/triggerhook/connection"
	"github.com/pvelx/triggerhook/contracts"
	"github.com/pvelx/triggerhook/error_service"
	"github.com/pvelx/triggerhook/monitoring_service"
	"path"
	"time"
)

func BuildTriggerHook(monitoring *Monitoring, conn connection.Options) contracts.TriggerHookInterface {

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
		contracts.PreloadingRate,
		contracts.CreatingRate,
		contracts.All,
		contracts.Preloaded,
		contracts.ConfirmationRate,
		contracts.WaitingForSending,
		contracts.WaitingForConfirmation,
		contracts.SendingRate,
		contracts.DeletingRate,
	} {
		topic := topic
		subscriptions[topic] = func(event contracts.MeasurementEvent) {
			monitoring.AddMeasurement(string(topic), event)
		}
	}

	tasksDeferredService := triggerhook.Build(triggerhook.Config{
		Connection: conn,
		ErrorServiceOptions: error_service.Options{
			Debug:         false,
			EventHandlers: eventHandlers,
		},
		MonitoringServiceOptions: monitoring_service.Options{
			PeriodMeasure: time.Second,
			Subscriptions: subscriptions,
		},
	})

	return tasksDeferredService
}
