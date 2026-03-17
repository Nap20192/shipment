package domain

import (
	"github.com/Nap20192/shipment/internal/pkg/kernel"
	"github.com/google/uuid"
)

type Shipment struct {
	ID            uuid.UUID
	Origin        string
	Destination   string
	Status        Status
	Cost          float64
	Revenue       float64
	Details       Details
	DriverDetails DriverDetails
	kernel.AggregateRoot
}

type Details struct {
	Weight     float64
	Dimensions [3]float64
}

type DriverDetails struct {
	Name string
}

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusCancelled Status = "CANCELLED"
)

const (
	BaseFee       float64 = 10.0
	RatePerKg     float64 = 2.0
	RatePerVolume float64 = 0.001
	RevenueRate   float64 = 0.80
)

func CalculateBasicCost(details Details) (cost float64, revenue float64) {
	volume := details.Dimensions[0] * details.Dimensions[1] * details.Dimensions[2]
	cost = BaseFee + (details.Weight * RatePerKg) + (volume * RatePerVolume)
	revenue = cost * RevenueRate
	return cost, revenue
}
