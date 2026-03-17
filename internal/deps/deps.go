package deps

import (
	"context"
	"fmt"

	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/core/domain/service"
	"github.com/Nap20192/shipment/internal/core/domain/spec"
	"github.com/Nap20192/shipment/internal/infra"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
	"github.com/Nap20192/shipment/internal/presentation/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	Repository app.ShipmentRepository
	AppService app.ShipmentService

	EventBus *app.EventBus
	Pool     *pgxpool.Pool
	Server   *grpc.Server
}
type DepsOption func(ctx context.Context, deps *Deps) error

func NewDeps(ctx context.Context, opts ...DepsOption) (*Deps, error) {
	deps := &Deps{}
	for _, opt := range opts {
		if err := opt(ctx, deps); err != nil {
			return nil, err
		}
	}
	return deps, nil
}

func WithRepository(connString string) DepsOption {
	return func(ctx context.Context, deps *Deps) error {
		config, err := pgxpool.ParseConfig(connString)

		if err != nil {
			return err
		}

		pool, err := pgxpool.NewWithConfig(ctx, config)

		if err != nil {
			return err
		}

		if err := pool.Ping(ctx); err != nil {
			return err
		}
		deps.Pool = pool
		deps.Repository = infra.NewRepo(sqlc.New(pool))
		return nil
	}
}

func WithGrpcServer(port string) DepsOption {
	return func(ctx context.Context, deps *Deps) error {
		server, err := grpc.NewServer(port, deps.AppService)
		if err != nil {
			return err
		}
		deps.Server = server
		return nil
	}
}

func WithShipmentService() DepsOption {
	return func(ctx context.Context, deps *Deps) error {

		rules, err := spec.DefaultTransitionSpec()
		if err != nil {
			return err
		}
		if deps.EventBus == nil || deps.Repository == nil {
			return fmt.Errorf("EventBus and Repository must be initialized before ShipmentService")
		}

		domainService := service.NewShipmentService(rules)
		service := app.NewShipmentService(domainService, deps.Repository, deps.EventBus)
		deps.AppService = service
		return nil
	}
}

func WithEventBus() DepsOption {
	return func(ctx context.Context, deps *Deps) error {
		deps.EventBus = app.NewEventBus()
		return nil
	}
}
