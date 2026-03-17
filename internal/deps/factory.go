package deps

import (
	"github.com/Nap20192/shipment/internal/core/app"
	domain_service "github.com/Nap20192/shipment/internal/core/domain/service"
	"github.com/Nap20192/shipment/internal/core/domain/spec"
	"github.com/Nap20192/shipment/internal/infra"
	"github.com/Nap20192/shipment/internal/pkg/sqlc"
)

func CreateAppService(queries *sqlc.Queries, eventBus *app.EventBus) (app.ShipmentService, error) {
	// 1. Create Repository Adapter
	repo := infra.NewSqlcShipmentRepository(queries)

	// 2. Create Domain Service
	statusSpec, err := spec.DefaultTransitionSpec()
	if err != nil {
		return nil, err
	}
	domainService := domain_service.NewShipmentService(statusSpec)

	// 3. Create Application Service
	return app.NewShipmentService(domainService, repo, eventBus), nil
}
