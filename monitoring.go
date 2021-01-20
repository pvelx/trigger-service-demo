package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/pvelx/triggerHook/contracts"
	"log"
	"time"
)

var sampleSize = 1000
var chPointCap = 10000

type Monitoring struct {
	connection client.Client
	chPoint    chan *client.Point
}

func NewMonitoring() *Monitoring {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "monitor",
		Password: "secret",
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	return &Monitoring{
		connection: c,
		chPoint:    make(chan *client.Point, chPointCap),
	}
}

func (m *Monitoring) Run() error {
	for {
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "trigger_hook",
			Precision: "ms",
		})
		if err != nil {
			return err
		}
		expire := time.After(5 * time.Second)
		for {
			select {
			case point := <-m.chPoint:
				bp.AddPoint(point)

				if len(bp.Points()) == sampleSize {
					goto done
				}

			case <-expire:
				goto done
			}
		}

	done:
		if len(bp.Points()) > 0 {
			err = m.connection.Write(bp)
			if err != nil {
				return err
			}
		}
	}
}

func (m *Monitoring) AddMeasurement(name string, event contracts.MeasurementEvent) {
	point, err := client.NewPoint(
		name,
		nil,
		map[string]interface{}{
			"value": event.Measurement,
		},
		event.Time,
	)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	m.chPoint <- point
}
