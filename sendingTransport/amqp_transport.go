package sendingTransport

import (
	"github.com/pvelx/triggerHook/contracts"
	"github.com/pvelx/triggerHook/domain"
)

func NewAmqpTransport() *amqpTransport {
	return &amqpTransport{}
}

type amqpTransport struct {
	contracts.SendingTransportInterface
}

func (amqpTransport) Send(task *domain.Task) bool {
	return true
}
