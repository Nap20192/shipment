package infra

import (
	"context"
	"fmt"

	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/pkg/kernel"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repo struct {
	queries *sqlc.Queries
}

func NewRepo(queries *sqlc.Queries) *Repo {
	return &Repo{
		queries: queries,
	}
}

func (r *Repo) Create(ctx context.Context, s domain.Shipment) error {
	_, err := r.queries.CreateShipment(ctx, sqlc.CreateShipmentParams{
		ID:              s.ID,
		Origin:          s.Origin,
		Destination:     s.Destination,
		Status:          string(s.Status),
		Cost:            floatToNumeric(s.Cost),
		Revenue:         floatToNumeric(s.Revenue),
		Weight:          floatToNumeric(s.Details.Weight),
		DimensionLength: floatToNumeric(s.Details.Dimensions[0]),
		DimensionWidth:  floatToNumeric(s.Details.Dimensions[1]),
		DimensionHeight: floatToNumeric(s.Details.Dimensions[2]),
		DriverName:      pgtype.Text{String: s.DriverDetails.Name, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create shipment: %w", err)
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (domain.Shipment, error) {
	dbShipment, err := r.queries.GetShipmentByID(ctx, id)
	if err != nil {
		return domain.Shipment{}, fmt.Errorf("failed to get shipment by id: %w", err)
	}

	return toDomainShipment(dbShipment), nil
}

func (r *Repo) AddEvent(ctx context.Context, shipmentID uuid.UUID, event kernel.DomainEvent) error {
	_, err := r.queries.AddShipmentEvent(ctx, sqlc.AddShipmentEventParams{
		ID:         uuid.New(),
		ShipmentID: shipmentID,
		EventName:  event.Name(),
		Payload:    event.Payload(),
	})
	if err != nil {
		return fmt.Errorf("failed to add shipment event: %w", err)
	}
	return nil
}

func (r *Repo) GetHistory(ctx context.Context, shipmentID uuid.UUID) ([]app.EventDTO, error) {
	events, err := r.queries.GetShipmentEventHistory(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment event history: %w", err)
	}

	dtos := make([]app.EventDTO, len(events))
	for i, e := range events {
		dtos[i] = app.EventDTO{
			ShipmentID: e.ShipmentID,
			EventName:  e.EventName,
			Payload:    e.Payload,
			CreatedAt:  e.CreatedAt.Time,
		}
	}
	return dtos, nil
}

func (r *Repo) UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, newStatus domain.Status) error {
	_, err := r.queries.UpdateShipmentStatus(ctx, sqlc.UpdateShipmentStatusParams{
		ID:     shipmentID,
		Status: string(newStatus),
	})
	if err != nil {
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	return nil
}

// Helpers

func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%f", f))
	return n
}

func numericToFloat(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

func toDomainShipment(s sqlc.Shipment) domain.Shipment {
	return domain.Shipment{
		ID:          s.ID,
		Origin:      s.Origin,
		Destination: s.Destination,
		Status:      domain.Status(s.Status),
		Cost:        numericToFloat(s.Cost),
		Revenue:     numericToFloat(s.Revenue),
		Details: domain.Details{
			Weight: numericToFloat(s.Weight),
			Dimensions: [3]float64{
				numericToFloat(s.DimensionLength),
				numericToFloat(s.DimensionWidth),
				numericToFloat(s.DimensionHeight),
			},
		},
		DriverDetails: domain.DriverDetails{
			Name: s.DriverName.String,
		},
	}
}
