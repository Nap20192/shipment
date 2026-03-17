package infra

import (
	"context"
	"fmt"

	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcShipmentRepository struct {
	queries *sqlc.Queries
}

func NewSqlcShipmentRepository(queries *sqlc.Queries) *SqlcShipmentRepository {
	return &SqlcShipmentRepository{
		queries: queries,
	}
}

func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(fmt.Sprintf("%f", f))
	return n
}

func numericToFloat(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

func (r *SqlcShipmentRepository) Create(ctx context.Context, s domain.Shipment) error {
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
	return err
}

func (r *SqlcShipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Shipment, error) {
	dbShipment, err := r.queries.GetShipmentByID(ctx, id)
	if err != nil {
		return domain.Shipment{}, err
	}

	return domain.Shipment{
		ID:          dbShipment.ID,
		Origin:      dbShipment.Origin,
		Destination: dbShipment.Destination,
		Status:      domain.Status(dbShipment.Status),
		Cost:        numericToFloat(dbShipment.Cost),
		Revenue:     numericToFloat(dbShipment.Revenue),
		Details: domain.Details{
			Weight: numericToFloat(dbShipment.Weight),
			Dimensions: [3]float64{
				numericToFloat(dbShipment.DimensionLength),
				numericToFloat(dbShipment.DimensionWidth),
				numericToFloat(dbShipment.DimensionHeight),
			},
		},
		DriverDetails: domain.DriverDetails{
			Name: dbShipment.DriverName.String,
		},
	}, nil
}

func (r *SqlcShipmentRepository) AddEvent(ctx context.Context, shipmentID uuid.UUID, status string, description string) error {
	_, err := r.queries.AddShipmentEvent(ctx, sqlc.AddShipmentEventParams{
		ID:          uuid.New(),
		ShipmentID:  shipmentID,
		Status:      status,
		Description: pgtype.Text{String: description, Valid: true},
	})
	return err
}

func (r *SqlcShipmentRepository) GetHistory(ctx context.Context, shipmentID uuid.UUID) ([]app.EventDTO, error) {
	events, err := r.queries.GetShipmentEventHistory(ctx, shipmentID)
	if err != nil {
		return nil, err
	}
	var dtos []app.EventDTO
	for _, e := range events {
		dtos = append(dtos, app.EventDTO{
			ShipmentID: e.ShipmentID,
			EventName:  e.Status,
			Payload:    []byte(e.Description.String),
			CreatedAt:  e.CreatedAt.Time,
		})
	}
	return dtos, nil
}
