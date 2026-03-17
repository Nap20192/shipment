package domain

import (
	"testing"
)

func TestCalculateBasicCost(t *testing.T) {
	tests := []struct {
		name        string
		details     Details
		wantCost    float64
		wantRevenue float64
	}{
		{
			name: "Basic calculation with positive dimensions and weight",
			details: Details{
				Weight:     10.0,
				Dimensions: [3]float64{10.0, 10.0, 10.0},
			},
			wantCost:    31.0,
			wantRevenue: 24.8,
		},
		{
			name: "Zero volume",
			details: Details{
				Weight:     5.0,
				Dimensions: [3]float64{0.0, 10.0, 10.0},
			},
			wantCost:    20.0,
			wantRevenue: 16.0,
		},
		{
			name: "Zero weight",
			details: Details{
				Weight:     0.0,
				Dimensions: [3]float64{10.0, 10.0, 10.0},
			},
			wantCost:    11.0,
			wantRevenue: 8.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCost, gotRevenue := CalculateBasicCost(tt.details)
			if gotCost != tt.wantCost {
				t.Errorf("CalculateBasicCost() gotCost = %v, want %v", gotCost, tt.wantCost)
			}
			if gotRevenue != tt.wantRevenue {
				t.Errorf("CalculateBasicCost() gotRevenue = %v, want %v", gotRevenue, tt.wantRevenue)
			}
		})
	}
}
