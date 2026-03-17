package app

import (
	"fmt"

	"github.com/Nap20192/shipment/internal/pkg/kernel"
)

type LogHandler struct {
}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) Handle(event kernel.DomainEvent) error {
	fmt.Printf("Event: %s, Payload: %+v\n", event.Name(), event.Payload())
	return nil
}
