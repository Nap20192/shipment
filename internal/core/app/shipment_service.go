package app

import "github.com/Nap20192/shipment/internal/pkg/sqlc"

type ShipmentService struct {
	shipmentRepo *sqlc.Queries
	mediator     *Mediator
}

func NewShipmentService(shipmentRepo *sqlc.Queries, mediator *Mediator) *ShipmentService {
	return &ShipmentService{
		shipmentRepo: shipmentRepo,
		mediator:     mediator,
	}
}
