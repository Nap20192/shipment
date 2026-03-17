package domain

import (
	"encoding/json"
)

const (
	ShipmentCreatedEventName       = "ShipmentCreated"
	ShipmentStatusUpdatedEventName = "ShipmentStatusUpdated"
)

type ShipmentCreatedEvent struct {
	ShipmentID  string `json:"shipment_id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

func (e ShipmentCreatedEvent) Name() string {
	return ShipmentCreatedEventName
}

func (e ShipmentCreatedEvent) Payload() []byte {
	payload, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return payload
}

type ShipmentStatusUpdatedEvent struct {
	ShipmentID string `json:"shipment_id"`
	NewStatus  Status `json:"new_status"`
}

func (e ShipmentStatusUpdatedEvent) Name() string {
	return ShipmentStatusUpdatedEventName
}

func (e ShipmentStatusUpdatedEvent) Payload() []byte {
	payload, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return payload
}
