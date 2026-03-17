package app

import (
	"context"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/pkg/kernel"
	"github.com/google/uuid"
)

type ShipmentRepository interface {
	Create(ctx context.Context, shipment domain.Shipment) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Shipment, error)
	AddEvent(ctx context.Context, shipmentID uuid.UUID, event kernel.DomainEvent) error
	GetHistory(ctx context.Context, shipmentID uuid.UUID) ([]EventDTO, error)
	UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, newStatus domain.Status) error
}
