package main

import (
	"github.com/pvelx/triggerHook/contracts"
	"github.com/pvelx/triggerHook/domain"
)

func NewTransportAmqp() *amqpTransport {
	return &amqpTransport{}
}

type amqpTransport struct {
	contracts.SendingTransportInterface
}

func (amqpTransport) Send(task *domain.Task) bool {
	return true
}
