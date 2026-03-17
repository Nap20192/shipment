package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Nap20192/shipment/internal/pkg/kernel"
)

type LogSubscriber struct {
	log *slog.Logger
}

func NewLogSubscriber() *LogSubscriber {
	return &LogSubscriber{
		log: slog.Default(),
	}
}

func (l *LogSubscriber) Handle(ctx context.Context, event kernel.DomainEvent) error {
	l.log.InfoContext(ctx, fmt.Sprintf("Event received: %s, payload: %+v", event.Name(), event.Payload()))
	return nil
}
