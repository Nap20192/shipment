package app

import "github.com/Nap20192/shipment/internal/pkg/kernel"

type EventHandler interface {
	Handle(event kernel.DomainEvent) error
}

type Mediator struct {
	handlers map[string][]EventHandler
}

func NewMediator() *Mediator {
	return &Mediator{
		handlers: make(map[string][]EventHandler),
	}
}

func (m *Mediator) Register(eventName string, handler EventHandler) {
	m.handlers[eventName] = append(m.handlers[eventName], handler)
}

func (m *Mediator) Publish(events []kernel.DomainEvent) error {
	for _, event := range events {
		if handlers, ok := m.handlers[event.Name()]; ok {
			for _, handler := range handlers {
				if err := handler.Handle(event); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
