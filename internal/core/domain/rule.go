package domain

import "fmt"

var ErrInvalidStatusTransition error = fmt.Errorf("invalid status transition")

type Rule interface {
	Check(shipment Shipment, newStatus Status) (bool, error)
}

type RuleFunc func(shipment Shipment, newStatus Status) (bool, error)

func (f RuleFunc) Check(shipment Shipment, newStatus Status) (bool, error) {
	return f(shipment, newStatus)
}

var AlwaysAllowRule Rule = RuleFunc(func(shipment Shipment, newStatus Status) (bool, error) {
	return true, nil
})

var AlwaysDenyRule Rule = RuleFunc(func(shipment Shipment, newStatus Status) (bool, error) {
	return false, ErrInvalidStatusTransition
})
