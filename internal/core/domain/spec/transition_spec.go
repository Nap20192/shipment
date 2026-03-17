package spec

import (
	"fmt"

	"github.com/Nap20192/shipment/internal/core/domain"
)

type StatusSpec interface {
	Check(shipment domain.Shipment, newStatus domain.Status) (bool, error)
}

type transitionValidation struct {
	registry map[domain.Status]map[domain.Status]domain.Rule
}

func (ts *transitionValidation) Check(shipment domain.Shipment, newStatus domain.Status) (bool, error) {
	if shipment.Status == newStatus {
		return false, fmt.Errorf("Error: %w, status: %s", domain.ErrAlreadyInStatus, newStatus)
	}
	if rules, ok := ts.registry[shipment.Status]; ok {
		if rule, ok := rules[newStatus]; ok {
			return rule.Check(shipment, newStatus)
		}
	}
	return false, domain.ErrInvalidStatusTransition
}

type TransitionValidationOption func(*transitionValidation)

func NewTransitionSpec(opts ...TransitionValidationOption) (*transitionValidation,error){

	ts := &transitionValidation{
		registry: make(map[domain.Status]map[domain.Status]domain.Rule),
	}

	for _, opt := range opts {
		opt(ts)
	}

	var missingEdges []string
	for _, from := range domain.AllStatuses {
		for _, to := range domain.AllStatuses {

			if from == to {
				continue
			}

			transitions, hasFrom := ts.registry[from]

			if !hasFrom {
				missingEdges = append(missingEdges, fmt.Sprintf("[%s -> %s]", from, to))
				continue
			}

			rule, hasRule := transitions[to]

			if !hasRule || rule == nil {
				missingEdges = append(missingEdges, fmt.Sprintf("[%s -> %s]", from, to))
			}
		}
	}

	if len(missingEdges) > 0 {
		return nil, fmt.Errorf("invalid adjacency matrix: missing transition rules for: %v", missingEdges)
	}

	return ts, nil
}

func WithRule(from domain.Status, to domain.Status, rule domain.Rule) TransitionValidationOption {
	return func(ts *transitionValidation) {
		if ts.registry[from] == nil {
			ts.registry[from] = make(map[domain.Status]domain.Rule)
		}
		ts.registry[from][to] = rule
	}
}

func DefaultTransitionSpec() (*transitionValidation, error) {
	return NewTransitionSpec(
		WithRule(domain.StatusPending, domain.StatusInTransit, domain.AlwaysAllowRule),
		WithRule(domain.StatusPending, domain.StatusCancelled, domain.AlwaysAllowRule),
		WithRule(domain.StatusPending, domain.StatusDelivered, domain.AlwaysDenyRule),

		WithRule(domain.StatusInTransit, domain.StatusDelivered, domain.AlwaysAllowRule),
		WithRule(domain.StatusInTransit, domain.StatusCancelled, domain.AlwaysAllowRule),
		WithRule(domain.StatusInTransit, domain.StatusPending, domain.AlwaysDenyRule),

		WithRule(domain.StatusDelivered, domain.StatusPending, domain.AlwaysDenyRule),
		WithRule(domain.StatusDelivered, domain.StatusInTransit, domain.AlwaysDenyRule),
		WithRule(domain.StatusDelivered, domain.StatusCancelled, domain.AlwaysDenyRule),

		WithRule(domain.StatusCancelled, domain.StatusPending, domain.AlwaysDenyRule),
		WithRule(domain.StatusCancelled, domain.StatusInTransit, domain.AlwaysDenyRule),
		WithRule(domain.StatusCancelled, domain.StatusDelivered, domain.AlwaysDenyRule),
	)
}
