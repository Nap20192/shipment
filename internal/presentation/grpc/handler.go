package grpc

import (
	"context"

	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/core/domain"
	pb "github.com/Nap20192/shipment/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedShipmentServiceServer
	service app.ShipmentService
}

func NewShipmentHandler(service app.ShipmentService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	if req.GetOrigin() == "" || req.GetDestination() == "" {
		return nil, status.Error(codes.InvalidArgument, "origin and destination are required")
	}

	details := domain.Details{
		Weight: req.GetDetails().GetWeight(),
		Dimensions: [3]float64{
			req.GetDetails().GetDimensionLength(),
			req.GetDetails().GetDimensionWidth(),
			req.GetDetails().GetDimensionHeight(),
		},
	}

	driverDetails := domain.DriverDetails{
		Name: req.GetDriverDetails().GetName(),
	}

	shipment, err := h.service.CreateShipment(ctx, req.GetOrigin(), req.GetDestination(), details, driverDetails)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create shipment: %v", err)
	}

	return &pb.CreateShipmentResponse{
		Shipment: domainShipmentToProto(shipment),
	}, nil
}

func (h *Handler) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	shipment, err := h.service.GetShipment(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shipment not found: %v", err)
	}

	return &pb.GetShipmentResponse{
		Shipment: domainShipmentToProto(shipment),
	}, nil
}

func (h *Handler) UpdateShipmentStatus(ctx context.Context, req *pb.UpdateShipmentStatusRequest) (*pb.UpdateShipmentStatusResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	newStatus, ok := protoStatusToDomain[req.GetNewStatus()]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid status")
	}

	shipment, err := h.service.UpdateShipmentStatus(ctx, req.GetId(), string(newStatus))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update shipment status: %v", err)
	}

	return &pb.UpdateShipmentStatusResponse{
		Shipment: domainShipmentToProto(shipment),
	}, nil
}

func (h *Handler) GetShipmentEventHistory(ctx context.Context, req *pb.GetShipmentEventHistoryRequest) (*pb.GetShipmentEventHistoryResponse, error) {
	if req.GetShipmentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "shipment_id is required")
	}

	events, err := h.service.History(ctx, req.GetShipmentId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get event history: %v", err)
	}

	pbEvents := make([]*pb.ShipmentEvent, 0, len(events))
	for _, e := range events {
		pbEvents = append(pbEvents, eventDTOToProto(e))
	}

	return &pb.GetShipmentEventHistoryResponse{
		Events: pbEvents,
	}, nil
}
