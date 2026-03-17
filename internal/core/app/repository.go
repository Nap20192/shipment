package app

import (
	"context"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/google/uuid"
)

type ShipmentRepository interface {
	Create(ctx context.Context, shipment domain.Shipment) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Shipment, error)
	AddEvent(ctx context.Context, shipmentID uuid.UUID, status string, description string) error
	GetHistory(ctx context.Context, shipmentID uuid.UUID) ([]EventDTO, error)
}
