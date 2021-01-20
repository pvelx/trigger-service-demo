module github.com/pvelx/triggerHookExample

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/influxdata/influxdb v1.8.3
	github.com/pvelx/triggerHook v0.0.0-00010101000000-000000000000
	github.com/streadway/amqp v1.0.0
	golang.org/x/net v0.0.0-20201201195509-5d6afe98e0b7 // indirect
	golang.org/x/sys v0.0.0-20201201145000-ef89a241ccb3 // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201201144952-b05cb90ed32e // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
)

replace github.com/pvelx/triggerHook => ../trigger-hook
