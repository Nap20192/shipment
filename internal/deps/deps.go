package deps

import (
	"context"

	"github.com/Nap20192/shipment/internal/core/app"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
)

type Deps struct {
	AppService app.ShipmentService
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

func WithShipmentService(queries *sqlc.Queries, eventBus *app.EventBus) DepsOption {
	return func(ctx context.Context, deps *Deps) error {
		service, err := CreateAppService(queries, eventBus)
		if err != nil {
			return err
		}
		deps.AppService = service
		return nil
	}
}
