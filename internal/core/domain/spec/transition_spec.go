package spec

import "github.com/Nap20192/shipment/internal/core/domain"


type StatusSpec interface {
	Check(shipment domain.Shipment, newStatus domain.Status) (bool, error)
}

type TransitionValidation struct {
	registry map[domain.Status]map[domain.Status]domain.Rule
}

func (ts *TransitionValidation) Check(shipment domain.Shipment, newStatus domain.Status) (bool, error) {
	if rules, ok := ts.registry[shipment.Status]; ok {
		if rule, ok := rules[newStatus]; ok {
			return rule.Check(shipment, newStatus)
		}
	}
	return false, domain.ErrInvalidStatusTransition
}

type TransitionValidationOption func(*TransitionValidation)

func NewTransitionSpec(opts ...TransitionValidationOption) *TransitionValidation {
	ts := &TransitionValidation{
		registry: make(map[domain.Status]map[domain.Status]domain.Rule),
	}
	for _, opt := range opts {
		opt(ts)
	}
	return ts
}

func WithRule(from domain.Status, to domain.Status, rule domain.Rule) TransitionValidationOption {
	return func(ts *TransitionValidation) {
		if ts.registry[from] == nil {
			ts.registry[from] = make(map[domain.Status]domain.Rule)
		}
		ts.registry[from][to] = rule
	}
}


var DefaultTransitionSpec = NewTransitionSpec(
	WithRule(domain.StatusPending, domain.StatusInTransit, domain.AlwaysAllowRule),
	WithRule(domain.StatusInTransit, domain.StatusDelivered, domain.AlwaysAllowRule),
	WithRule(domain.StatusPending, domain.StatusCancelled, domain.AlwaysAllowRule),
	WithRule(domain.StatusInTransit, domain.StatusCancelled, domain.AlwaysAllowRule),
)
