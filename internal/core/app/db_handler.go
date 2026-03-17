package app

import (
	"github.com/Nap20192/shipment/internal/pkg/kernel"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
)

type DBHandler struct {
	queries *sqlc.Queries
}

func NewDBHandler(queries *sqlc.Queries) *DBHandler {
	return &DBHandler{
		queries: queries,
	}
}

func (h *DBHandler) Handle(event kernel.DomainEvent) error {
	return nil
}
