package app

import (
	"fmt"

	"github.com/Nap20192/shipment/internal/pkg/kernel"
)

type LogSubscriber struct {
}

func NewLogSubscriber() *LogSubscriber {
	return &LogSubscriber{}
}

func (l *LogSubscriber) Handle(event kernel.DomainEvent) error {
	fmt.Printf("Event: %s, Payload: %+v\n", event.Name(), event.Payload())
	return nil
}
