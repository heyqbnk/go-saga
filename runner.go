package gotransaction

import (
	"context"
)

type Runner struct {
	resultStack []func(ctx context.Context) error
}

func newRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Rollback(ctx context.Context) error {
	for i := len(r.resultStack) - 1; i >= 0; i-- {
		if err := r.resultStack[i](ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) runStep(ctx context.Context, step Step[any]) (any, error) {
	res, err := retry[any](func() (any, error) {
		return step.Run(ctx)
	}, step.ShouldRetryRun)
	if err != nil {
		return nil, err
	}

	// Cache result.
	if step.Rollback != nil {
		r.resultStack = append(r.resultStack, func(ctx context.Context) error {
			return retryWithoutResult(func() error {
				return step.Rollback(ctx, res)
			}, step.ShouldRetryRollback)
		})
	}

	return res, err
}
