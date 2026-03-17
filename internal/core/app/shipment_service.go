package app

import (
	"context"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
)

type ShipmentService interface {
	UpdateShipmentStatus(context.Context, string, string) (sqlc.Shipment, error)
	CreateShipment(context.Context, string, string, domain.Details, domain.DriverDetails) (domain.Shipment, error)
}
type shipmentService struct {
	shipmentRepo *sqlc.Queries
	EventBus	 *EventBus
}

func NewShipmentService(shipmentRepo *sqlc.Queries, eventBus *EventBus) *shipmentService {
	return &shipmentService{
		shipmentRepo: shipmentRepo,
		EventBus:     eventBus,
	}
}
