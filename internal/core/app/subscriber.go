package app

import (
	"context"
	"fmt"

	"github.com/Nap20192/shipment/internal/pkg/kernel"
)

type LogSubscriber struct {
}

func NewLogSubscriber() *LogSubscriber {
	return &LogSubscriber{}
}

func (l *LogSubscriber) Handle(ctx context.Context, event kernel.DomainEvent) error {
	fmt.Printf("Event received: %s, Payload: %s\n", event.Name(), event.Payload())
	return nil
}
