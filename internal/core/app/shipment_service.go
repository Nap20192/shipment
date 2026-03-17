package app

import (
	"context"
	"time"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/core/domain/service"
	"github.com/google/uuid"
)

var EventBusKey string = "shipment_service"

type EventDTO struct {
	ShipmentID uuid.UUID
	EventName  string
	Payload    []byte
	CreatedAt  time.Time
}

type ShipmentService interface {
	UpdateShipmentStatus(context.Context, string, string) (domain.Shipment, error)
	CreateShipment(context.Context, string, string, domain.Details, domain.DriverDetails) (domain.Shipment, error)
	History(context.Context, string) ([]EventDTO, error)
	GetShipment(context.Context, string) (domain.Shipment, error)
}

type shipmentService struct {
	domainService service.ShipmentService
	repo          ShipmentRepository
	EventBus      *EventBus
}

func NewShipmentService(domainService service.ShipmentService, repo ShipmentRepository, eventBus *EventBus) *shipmentService {
	return &shipmentService{
		domainService: domainService,
		repo:          repo,
		EventBus:      eventBus,
	}
}

func (s *shipmentService) CreateShipment(ctx context.Context, origin, destination string, details domain.Details, driverDetails domain.DriverDetails) (domain.Shipment, error) {
	shipment, err := s.domainService.CreateShipment(origin, destination, details, driverDetails)
	if err != nil {
		return domain.Shipment{}, err
	}

	err = s.repo.Create(ctx, shipment)
	if err != nil {
		return domain.Shipment{}, err
	}
	events := shipment.DomainEvents()

	for _, event := range events {
		err = s.repo.AddEvent(ctx, shipment.ID, event)
		if err != nil {
			return domain.Shipment{}, err
		}
		err = s.EventBus.Publish(ctx, EventBusKey, event)
		if err != nil {
			return domain.Shipment{}, err
		}
	}

	return shipment, nil
}

func (s *shipmentService) UpdateShipmentStatus(ctx context.Context, id string, statusStr string) (domain.Shipment, error) {
	shipment, err := s.GetShipment(ctx, id)
	if err != nil {
		return domain.Shipment{}, err
	}
	shipment, err = s.domainService.UpdateShipmentStatus(shipment, domain.Status(statusStr))
	if err != nil {
		return domain.Shipment{}, err
	}
	err = s.repo.UpdateShipmentStatus(ctx, shipment.ID, shipment.Status)
	if err != nil {
		return domain.Shipment{}, err
	}
	events := shipment.DomainEvents()
	for _, event := range events {
		err = s.repo.AddEvent(ctx, shipment.ID, event)
		if err != nil {
			return domain.Shipment{}, err
		}
		err = s.EventBus.Publish(ctx, EventBusKey, event)
		if err != nil {
			return domain.Shipment{}, err
		}
	}

	return shipment, nil
}

func (s *shipmentService) GetShipment(ctx context.Context, id string) (domain.Shipment, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return domain.Shipment{}, err
	}
	return s.repo.GetByID(ctx, parsedID)
}

func (s *shipmentService) History(ctx context.Context, id string) ([]EventDTO, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetHistory(ctx, parsedID)
}
