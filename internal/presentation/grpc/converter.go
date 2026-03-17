package grpc

import (
	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/core/domain"
	pb "github.com/Nap20192/shipment/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainStatusToProto = map[domain.Status]pb.ShipmentStatus{
	domain.StatusPending:   pb.ShipmentStatus_SHIPMENT_STATUS_PENDING,
	domain.StatusInTransit: pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT,
	domain.StatusDelivered: pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED,
	domain.StatusCancelled: pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED,
}

var protoStatusToDomain = map[pb.ShipmentStatus]domain.Status{
	pb.ShipmentStatus_SHIPMENT_STATUS_PENDING:    domain.StatusPending,
	pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT: domain.StatusInTransit,
	pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:  domain.StatusDelivered,
	pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED:  domain.StatusCancelled,
}

func domainShipmentToProto(s domain.Shipment) *pb.Shipment {
	return &pb.Shipment{
		Id:          s.ID.String(),
		Origin:      s.Origin,
		Destination: s.Destination,
		Status:      domainStatusToProto[s.Status],
		Cost:        s.Cost,
		Revenue:     s.Revenue,
		Details: &pb.ShipmentDetails{
			Weight:          s.Details.Weight,
			DimensionLength: s.Details.Dimensions[0],
			DimensionWidth:  s.Details.Dimensions[1],
			DimensionHeight: s.Details.Dimensions[2],
		},
		DriverDetails: &pb.DriverDetails{
			Name: s.DriverDetails.Name,
		},
	}
}

func eventDTOToProto(e app.EventDTO) *pb.ShipmentEvent {
	ev := &pb.ShipmentEvent{
		EventName: e.EventName,
		Payload:   e.Payload,
	}
	if !e.CreatedAt.IsZero() {
		ev.CreatedAt = timestamppb.New(e.CreatedAt)
	}
	return ev
}
