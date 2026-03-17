package service

import (
	"fmt"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/core/domain/spec"
	"github.com/google/uuid"
)
type ShipmentService interface {
	UpdateShipmentStatus(shipment domain.Shipment, newStatus domain.Status) (domain.Shipment, error)
	CreateShipment(origin, destination string, details domain.Details, driverDetails domain.DriverDetails) (domain.Shipment, error)
}

var (
	ErrInvalidShipment error = fmt.Errorf("invalid shipment details")
)

type shipmentService struct {
	statusSpec spec.StatusSpec
}

func NewShipmentService(statusSpec spec.StatusSpec) *shipmentService {
	return &shipmentService{
		statusSpec: statusSpec,
	}
}

func (s *shipmentService) UpdateShipmentStatus(shipment domain.Shipment, newStatus domain.Status) (domain.Shipment, error) {
	allowed, err := s.statusSpec.Check(shipment, newStatus)
	if err != nil {
		return shipment, err
	}
	if !allowed {
		return shipment, domain.ErrInvalidStatusTransition
	}

	shipment.Status = newStatus

	shipment.ApplyDomain(domain.ShipmentStatusUpdatedEvent{
		ShipmentID: shipment.ID.String(),
		NewStatus:  newStatus,
	})

	return shipment, nil
}

func (s *shipmentService) CreateShipment(origin, destination string, details domain.Details, driverDetails domain.DriverDetails) (domain.Shipment, error) {
	if origin == "" || destination == "" {
		return domain.Shipment{}, ErrInvalidShipment
	}
	if details.Weight <= 0 || details.Dimensions[0] <= 0 || details.Dimensions[1] <= 0 || details.Dimensions[2] <= 0 {
		return domain.Shipment{}, ErrInvalidShipment
	}

	cost, revenue := domain.CalculateBasicCost(details)
	shipment := domain.Shipment{
		ID:            uuid.New(),
		Origin:        origin,
		Destination:   destination,
		Status:        domain.StatusPending,
		Cost:          cost,
		Revenue:       revenue,
		Details:       details,
		DriverDetails: driverDetails,
	}
	shipment.ApplyDomain(domain.ShipmentCreatedEvent{
		ShipmentID:  shipment.ID.String(),
		Origin:      shipment.Origin,
		Destination: shipment.Destination,
	})

	return shipment, nil
}
