package infra

import (
	"context"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/pkg/kernel"
)

type Repository interface {
	Create(ctx context.Context, shipment domain.Shipment) error
	Update(ctx context.Context, id string, newStatus string) error

	GetById(ctx context.Context, id string) (domain.Shipment, error)
	History(ctx context.Context, id string) ([]kernel.DomainEvent, error)
}
