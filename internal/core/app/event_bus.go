package app

import (
	"context"
	"errors"
	"sync"

	"github.com/Nap20192/shipment/internal/pkg/kernel"
)

type EventHandler interface {
	Handle(ctx context.Context, event kernel.DomainEvent) error
}

type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (b *EventBus) Subscribe(key string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[key] = append(b.handlers[key], handler)
}

func (b *EventBus) Publish(ctx context.Context, key string, event kernel.DomainEvent) error {
	var errs []error
	b.mu.RLock()
	handlers := b.handlers[key]
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
