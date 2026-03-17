package grpc

import (
	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
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

func sqlcShipmentToProto(s sqlc.Shipment) *pb.Shipment {
	out := &pb.Shipment{
		Id:          s.ID.String(),
		Origin:      s.Origin,
		Destination: s.Destination,
	}

	if st, ok := domainStatusToProto[domain.Status(s.Status)]; ok {
		out.Status = st
	}

	if f, err := s.Cost.Float64Value(); err == nil {
		out.Cost = f.Float64
	}
	if f, err := s.Revenue.Float64Value(); err == nil {
		out.Revenue = f.Float64
	}

	details := &pb.ShipmentDetails{}
	if f, err := s.Weight.Float64Value(); err == nil {
		details.Weight = f.Float64
	}
	if f, err := s.DimensionLength.Float64Value(); err == nil {
		details.DimensionLength = f.Float64
	}
	if f, err := s.DimensionWidth.Float64Value(); err == nil {
		details.DimensionWidth = f.Float64
	}
	if f, err := s.DimensionHeight.Float64Value(); err == nil {
		details.DimensionHeight = f.Float64
	}
	out.Details = details

	out.DriverDetails = &pb.DriverDetails{
		Name: s.DriverName.String,
	}

	if s.CreatedAt.Valid {
		out.CreatedAt = timestamppb.New(s.CreatedAt.Time)
	}
	if s.UpdatedAt.Valid {
		out.UpdatedAt = timestamppb.New(s.UpdatedAt.Time)
	}

	return out
}

func sqlcEventToProto(e sqlc.ShipmentEvent) *pb.ShipmentEvent {
	out := &pb.ShipmentEvent{
		Id:         e.ID.String(),
		ShipmentId: e.ShipmentID.String(),
		EventName:  e.Status,
		Payload:    []byte(e.Description.String),
	}

	if e.CreatedAt.Valid {
		out.CreatedAt = timestamppb.New(e.CreatedAt.Time)
	}

	return out
}
