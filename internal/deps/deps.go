package deps

import "context"

type Deps struct {
	// Add any shared dependencies here, e.g. database connections, configuration, etc.
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
