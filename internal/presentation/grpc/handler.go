package grpc

import (
	"context"

	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/core/domain"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
	pb "github.com/Nap20192/shipment/proto/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShipmentHandler struct {
	pb.UnimplementedShipmentServiceServer
	service app.ShipmentService
	queries *sqlc.Queries
}

func NewShipmentHandler(service app.ShipmentService, queries *sqlc.Queries) *ShipmentHandler {
	return &ShipmentHandler{
		service: service,
		queries: queries,
	}
}

func (h *ShipmentHandler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
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

func (h *ShipmentHandler) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shipment id")
	}

	shipment, err := h.queries.GetShipmentByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shipment not found: %v", err)
	}

	return &pb.GetShipmentResponse{
		Shipment: sqlcShipmentToProto(shipment),
	}, nil
}

func (h *ShipmentHandler) UpdateShipmentStatus(ctx context.Context, req *pb.UpdateShipmentStatusRequest) (*pb.UpdateShipmentStatusResponse, error) {
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
		Shipment: sqlcShipmentToProto(shipment),
	}, nil
}

func (h *ShipmentHandler) GetShipmentEventHistory(ctx context.Context, req *pb.GetShipmentEventHistoryRequest) (*pb.GetShipmentEventHistoryResponse, error) {
	if req.GetShipmentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "shipment_id is required")
	}

	id, err := uuid.Parse(req.GetShipmentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid shipment id")
	}

	events, err := h.queries.GetShipmentEventHistory(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get event history: %v", err)
	}

	pbEvents := make([]*pb.ShipmentEvent, 0, len(events))
	for _, e := range events {
		pbEvents = append(pbEvents, sqlcEventToProto(e))
	}

	return &pb.GetShipmentEventHistoryResponse{
		Events: pbEvents,
	}, nil
}
