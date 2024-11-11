package databasse

import (
	"context"
)

type Chain struct {
	healths []Health
}

func NewChain(healths ...Health) Chain {
	return Chain{healths: healths}
}

func (h Chain) Readiness(ctx context.Context) error {
	for _, check := range h.healths {
		if err := check.Readiness(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (h Chain) Liveness(ctx context.Context) error {
	for _, check := range h.healths {
		if err := check.Liveness(ctx); err != nil {
			return err
		}
	}

	return nil
}
