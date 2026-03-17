package service

import (
	"errors"
	"testing"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/google/uuid"
)

type mockStatusSpec struct {
	allow bool
	err   error
}

func (m *mockStatusSpec) Check(shipment domain.Shipment, newStatus domain.Status) (bool, error) {
	return m.allow, m.err
}

func TestCreateShipment(t *testing.T) {
	s := NewShipmentService(&mockStatusSpec{})

	tests := []struct {
		name          string
		origin        string
		destination   string
		details       domain.Details
		driverDetails domain.DriverDetails
		wantErr       error
	}{
		{
			name:        "Successful creation",
			origin:      "A",
			destination: "B",
			details: domain.Details{
				Weight:     10,
				Dimensions: [3]float64{1, 1, 1},
			},
			driverDetails: domain.DriverDetails{Name: "Driver"},
			wantErr:       nil,
		},
		{
			name:        "Empty origin",
			origin:      "",
			destination: "B",
			details: domain.Details{
				Weight:     10,
				Dimensions: [3]float64{1, 1, 1},
			},
			driverDetails: domain.DriverDetails{Name: "Driver"},
			wantErr:       ErrInvalidShipment,
		},
		{
			name:        "Zero weight",
			origin:      "A",
			destination: "B",
			details: domain.Details{
				Weight:     0,
				Dimensions: [3]float64{1, 1, 1},
			},
			driverDetails: domain.DriverDetails{Name: "Driver"},
			wantErr:       ErrInvalidShipment,
		},
		{
			name:        "Negative weight",
			origin:      "A",
			destination: "B",
			details: domain.Details{
				Weight:     -1,
				Dimensions: [3]float64{1, 1, 1},
			},
			driverDetails: domain.DriverDetails{Name: "Driver"},
			wantErr:       ErrInvalidShipment,
		},
		{
			name:        "Zero dimension",
			origin:      "A",
			destination: "B",
			details: domain.Details{
				Weight:     10,
				Dimensions: [3]float64{1, 0, 1},
			},
			driverDetails: domain.DriverDetails{Name: "Driver"},
			wantErr:       ErrInvalidShipment,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreateShipment(tt.origin, tt.destination, tt.details, tt.driverDetails)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateShipment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if got.ID == uuid.Nil {
					t.Errorf("CreateShipment() got ID = %v, want non-nil", got.ID)
				}
				if got.Origin != tt.origin {
					t.Errorf("CreateShipment() got Origin = %v, want %v", got.Origin, tt.origin)
				}
				if len(got.DomainEvents()) != 1 {
					t.Errorf("CreateShipment() got %v events, want 1", len(got.DomainEvents()))
				}
			}
		})
	}
}

func TestUpdateShipmentStatus(t *testing.T) {
	tests := []struct {
		name      string
		shipment  domain.Shipment
		newStatus domain.Status
		mockAllow bool
		mockErr   error
		wantErr   error
	}{
		{
			name:      "Allowed transition",
			shipment:  domain.Shipment{ID: uuid.New(), Status: domain.StatusPending},
			newStatus: domain.StatusInTransit,
			mockAllow: true,
			mockErr:   nil,
			wantErr:   nil,
		},
		{
			name:      "Denied by spec error",
			shipment:  domain.Shipment{ID: uuid.New(), Status: domain.StatusPending},
			newStatus: domain.StatusDelivered,
			mockAllow: false,
			mockErr:   domain.ErrInvalidStatusTransition,
			wantErr:   domain.ErrInvalidStatusTransition,
		},
		{
			name:      "Denied by spec allow=false",
			shipment:  domain.Shipment{ID: uuid.New(), Status: domain.StatusPending},
			newStatus: domain.StatusDelivered,
			mockAllow: false,
			mockErr:   nil,
			wantErr:   domain.ErrInvalidStatusTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShipmentService(&mockStatusSpec{allow: tt.mockAllow, err: tt.mockErr})
			got, err := s.UpdateShipmentStatus(tt.shipment, tt.newStatus)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateShipmentStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if got.Status != tt.newStatus {
					t.Errorf("UpdateShipmentStatus() got status = %v, want %v", got.Status, tt.newStatus)
				}
				if len(got.DomainEvents()) != 1 {
					t.Errorf("UpdateShipmentStatus() got %v events, want 1", len(got.DomainEvents()))
				}
			}
		})
	}
}
