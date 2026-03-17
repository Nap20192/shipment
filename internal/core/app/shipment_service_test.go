package app

import (
	"context"
	"errors"
	"testing"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/pkg/kernel"
	"github.com/google/uuid"
)

// Mocks

type mockDomainService struct {
	createFn func(origin, destination string, details domain.Details, driverDetails domain.DriverDetails) (domain.Shipment, error)
	updateFn func(shipment domain.Shipment, newStatus domain.Status) (domain.Shipment, error)
}

func (m *mockDomainService) CreateShipment(origin, destination string, details domain.Details, driverDetails domain.DriverDetails) (domain.Shipment, error) {
	return m.createFn(origin, destination, details, driverDetails)
}

func (m *mockDomainService) UpdateShipmentStatus(shipment domain.Shipment, newStatus domain.Status) (domain.Shipment, error) {
	return m.updateFn(shipment, newStatus)
}

type mockRepo struct {
	ShipmentRepository
	createFn       func(ctx context.Context, shipment domain.Shipment) error
	getByIDFn      func(ctx context.Context, id uuid.UUID) (domain.Shipment, error)
	addEventFn     func(ctx context.Context, shipmentID uuid.UUID, event kernel.DomainEvent) error
	updateStatusFn func(ctx context.Context, shipmentID uuid.UUID, newStatus domain.Status) error
}

func (m *mockRepo) Create(ctx context.Context, s domain.Shipment) error { return m.createFn(ctx, s) }
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Shipment, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockRepo) AddEvent(ctx context.Context, id uuid.UUID, e kernel.DomainEvent) error {
	return m.addEventFn(ctx, id, e)
}

func (m *mockRepo) UpdateShipmentStatus(ctx context.Context, id uuid.UUID, s domain.Status) error {
	return m.updateStatusFn(ctx, id, s)
}

func TestAppCreateShipment(t *testing.T) {
	ctx := context.Background()
	eb := NewEventBus()

	shipmentID := uuid.New()

	mockDS := &mockDomainService{
		createFn: func(origin, destination string, details domain.Details, driverDetails domain.DriverDetails) (domain.Shipment, error) {
			s := domain.Shipment{ID: shipmentID, Origin: origin, Destination: destination}
			s.ApplyDomain(domain.ShipmentCreatedEvent{ShipmentID: shipmentID.String()})
			return s, nil
		},
	}

	repoCalled := false
	eventAddedCalled := false
	mockR := &mockRepo{
		createFn: func(ctx context.Context, shipment domain.Shipment) error {
			repoCalled = true
			return nil
		},
		addEventFn: func(ctx context.Context, shipmentID uuid.UUID, event kernel.DomainEvent) error {
			eventAddedCalled = true
			return nil
		},
	}

	s := NewShipmentService(mockDS, mockR, eb)

	_, err := s.CreateShipment(ctx, "A", "B", domain.Details{Weight: 1}, domain.DriverDetails{Name: "D"})
	if err != nil {
		t.Fatalf("CreateShipment failed: %v", err)
	}

	if !repoCalled {
		t.Error("Repository.Create was not called")
	}
	if !eventAddedCalled {
		t.Error("Repository.AddEvent was not called")
	}
}

func TestAppUpdateShipmentStatus(t *testing.T) {
	ctx := context.Background()
	eb := NewEventBus()
	shipmentID := uuid.New()

	mockDS := &mockDomainService{
		updateFn: func(s domain.Shipment, newStatus domain.Status) (domain.Shipment, error) {
			s.Status = newStatus
			s.ApplyDomain(domain.ShipmentStatusUpdatedEvent{ShipmentID: s.ID.String(), NewStatus: newStatus})
			return s, nil
		},
	}

	mockR := &mockRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (domain.Shipment, error) {
			return domain.Shipment{ID: shipmentID, Status: domain.StatusPending}, nil
		},
		updateStatusFn: func(ctx context.Context, id uuid.UUID, s domain.Status) error {
			return nil
		},
		addEventFn: func(ctx context.Context, id uuid.UUID, e kernel.DomainEvent) error {
			return nil
		},
	}

	s := NewShipmentService(mockDS, mockR, eb)

	_, err := s.UpdateShipmentStatus(ctx, shipmentID.String(), "IN_TRANSIT")
	if err != nil {
		t.Fatalf("UpdateShipmentStatus failed: %v", err)
	}
}

func TestAppUpdateShipmentStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	mockR := &mockRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (domain.Shipment, error) {
			return domain.Shipment{}, errors.New("not found")
		},
	}
	s := NewShipmentService(&mockDomainService{}, mockR, NewEventBus())

	_, err := s.UpdateShipmentStatus(ctx, uuid.New().String(), "IN_TRANSIT")
	if err == nil {
		t.Error("Expected error for non-existent shipment, got nil")
	}
}
