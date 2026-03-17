package domain

import (
	"encoding/json"
)

type ShipmentCreatedEvent struct {
	ShipmentID  string `json:"shipment_id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

func (e ShipmentCreatedEvent) Name() string {
	return "ShipmentCreated"
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
	NewStatus  string `json:"new_status"`
}

func (e ShipmentStatusUpdatedEvent) Name() string {
	return "ShipmentStatusUpdated"
}

func (e ShipmentStatusUpdatedEvent) Payload() []byte {
	payload, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return payload
}
