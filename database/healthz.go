package databasse

import "context"

type Health interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}

type health struct {
	readinessFunc func(ctx context.Context) error
	livenessFunc  func(ctx context.Context) error
}

func New(readiness, liveness func(ctx context.Context) error) Health {
	return health{
		readinessFunc: readiness,
		livenessFunc:  liveness,
	}
}

func (h health) Readiness(ctx context.Context) error {
	return h.readinessFunc(ctx)
}

func (h health) Liveness(ctx context.Context) error {
	return h.livenessFunc(ctx)
}
