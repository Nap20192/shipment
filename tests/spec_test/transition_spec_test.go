package spec_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/core/domain/spec"
)

var errMockRuleDeny = errors.New("mock rule denied transition")

type mockRule struct {
	allow bool
	err   error
}

func (m mockRule) Check(shipment domain.Shipment, newStatus domain.Status) (bool, error) {
	return m.allow, m.err
}

func TestTransitionValidation_Check(t *testing.T) {
	ts, err := spec.NewTransitionSpec(
		spec.WithRule(domain.StatusPending, domain.StatusInTransit, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusPending, domain.StatusDelivered, mockRule{allow: false, err: errMockRuleDeny}),
		spec.WithRule(domain.StatusPending, domain.StatusCancelled, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusInTransit, domain.StatusPending, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusInTransit, domain.StatusDelivered, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusInTransit, domain.StatusCancelled, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusDelivered, domain.StatusPending, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusDelivered, domain.StatusInTransit, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusDelivered, domain.StatusCancelled, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusCancelled, domain.StatusPending, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusCancelled, domain.StatusInTransit, mockRule{allow: true, err: nil}),
		spec.WithRule(domain.StatusCancelled, domain.StatusDelivered, mockRule{allow: true, err: nil}),
	)

	if err != nil {
		t.Fatalf("Failed to setup test spec: %v", err)
	}

	tests := []struct {
		name      string
		shipment  domain.Shipment
		newStatus domain.Status
		wantAllow bool
		wantErr   error
	}{
		{
			name:      "Self transition should fail with ErrAlreadyInStatus",
			shipment:  domain.Shipment{Status: domain.StatusPending},
			newStatus: domain.StatusPending,
			wantAllow: false,
			wantErr:   domain.ErrAlreadyInStatus,
		},
		{
			name:      "Valid transition allowed by rule",
			shipment:  domain.Shipment{Status: domain.StatusPending},
			newStatus: domain.StatusInTransit,
			wantAllow: true,
			wantErr:   nil,
		},
		{
			name:      "Valid transition path but denied by rule logic",
			shipment:  domain.Shipment{Status: domain.StatusPending},
			newStatus: domain.StatusDelivered,
			wantAllow: false,
			wantErr:   errMockRuleDeny,
		},
		{
			name:      "Unregistered transition path should return ErrInvalidStatusTransition",
			shipment:  domain.Shipment{Status: domain.StatusPending},
			newStatus: domain.Status("UNKNOWN_FAKE_STATUS"),
			wantAllow: false,
			wantErr:   domain.ErrInvalidStatusTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAllow, gotErr := ts.Check(tt.shipment, tt.newStatus)

			if gotAllow != tt.wantAllow {
				t.Errorf("Check() gotAllow = %v, want %v", gotAllow, tt.wantAllow)
			}
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Check() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestNewTransitionSpec(t *testing.T) {
	t.Run("Successfully creates a fully populated matrix", func(t *testing.T) {
		var opts []spec.TransitionValidationOption

		for _, from := range domain.AllStatuses {
			for _, to := range domain.AllStatuses {
				if from == to {
					continue
				}
				opts = append(opts, spec.WithRule(from, to, mockRule{allow: true}))
			}
		}

		ts, err := spec.NewTransitionSpec(opts...)
		if err != nil {
			t.Fatalf("Expected successful creation, got error: %v", err)
		}
		if ts == nil {
			t.Fatal("Expected a non-nil spec object")
		}
	})

	t.Run("Fails to create when edges are missing", func(t *testing.T) {
		ts, err := spec.NewTransitionSpec()

		if err == nil {
			t.Fatal("Expected an error due to empty adjacency matrix, got nil")
		}
		if ts != nil {
			t.Fatal("Spec should be nil when matrix validation fails")
		}

		if !strings.Contains(err.Error(), "missing transition rules for:") {
			t.Errorf("Unexpected error format: %v", err)
		}
	})
}

func TestDefaultTransitionSpec(t *testing.T) {
	t.Run("Default configuration covers all required edges", func(t *testing.T) {
		ts, err := spec.DefaultTransitionSpec()

		if err != nil {
			t.Fatalf("DefaultTransitionSpec is missing edges and failed to initialize: %v\nCheck if domain.AllStatuses matches the rules inside DefaultTransitionSpec.", err)
		}

		if ts == nil {
			t.Fatal("DefaultTransitionSpec returned a nil spec")
		}
	})

	t.Run("Default configuration rules behave as expected", func(t *testing.T) {
		ts, err := spec.DefaultTransitionSpec()
		if err != nil {
			t.Fatalf("Failed to initialize default spec: %v", err)
		}

		tests := []struct {
			from      domain.Status
			to        domain.Status
			wantAllow bool
		}{
			{from: domain.StatusPending, to: domain.StatusInTransit, wantAllow: true},
			{from: domain.StatusPending, to: domain.StatusCancelled, wantAllow: true},
			{from: domain.StatusPending, to: domain.StatusDelivered, wantAllow: false},

			{from: domain.StatusInTransit, to: domain.StatusDelivered, wantAllow: true},
			{from: domain.StatusInTransit, to: domain.StatusCancelled, wantAllow: true},
			{from: domain.StatusInTransit, to: domain.StatusPending, wantAllow: false},

			{from: domain.StatusDelivered, to: domain.StatusPending, wantAllow: false},
			{from: domain.StatusDelivered, to: domain.StatusInTransit, wantAllow: false},
			{from: domain.StatusDelivered, to: domain.StatusCancelled, wantAllow: false},

			{from: domain.StatusCancelled, to: domain.StatusPending, wantAllow: false},
			{from: domain.StatusCancelled, to: domain.StatusInTransit, wantAllow: false},
			{from: domain.StatusCancelled, to: domain.StatusDelivered, wantAllow: false},
		}

		for _, tt := range tests {
			testName := fmt.Sprintf("%s -> %s", tt.from, tt.to)

			t.Run(testName, func(t *testing.T) {
				shipment := domain.Shipment{Status: tt.from}
				allow, err := ts.Check(shipment, tt.to)

				if allow != tt.wantAllow {
					t.Errorf("Check() allow = %v, want %v", allow, tt.wantAllow)
				}

				if tt.wantAllow && err != nil {
					t.Errorf("Expected allowed transition to have no error, but got: %v", err)
				}
				if !tt.wantAllow && err == nil {
					t.Errorf("Expected denied transition to return an error, but got nil")
				}
			})
		}
	})
}

type dummySpec struct {
	registry map[domain.Status]map[domain.Status]domain.Rule
}

func (d *dummySpec) Check(shipment domain.Shipment, newStatus domain.Status) (bool, error) {
	if shipment.Status == newStatus {
		return false, domain.ErrAlreadyInStatus
	}

	if rules, ok := d.registry[shipment.Status]; ok {
		if rule, ok := rules[newStatus]; ok {
			return rule.Check(shipment, newStatus)
		}
	}

	return false, domain.ErrInvalidStatusTransition
}

func NewTransitionSpecForTestOnly() spec.StatusSpec {
	return &dummySpec{
		registry: make(map[domain.Status]map[domain.Status]domain.Rule),
	}
}
