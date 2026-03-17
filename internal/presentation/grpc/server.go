package grpc

import (
	"net"

	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
	pb "github.com/Nap20192/shipment/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(addr string, service app.ShipmentService, queries *sqlc.Queries, opts ...grpc.ServerOption) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(opts...)

	handler := NewShipmentHandler(service, queries)
	pb.RegisterShipmentServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		listener:   lis,
	}, nil
}

func (s *Server) Serve() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}
