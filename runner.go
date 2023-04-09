package saga

import (
	"context"
	"fmt"
)

type Runner struct {
	resultStack []func(ctx context.Context) error
}

func newRunner() *Runner {
	return &Runner{}
}

// Rollback rolls back previously completed steps in reverse order.
func (r *Runner) Rollback(ctx context.Context) error {
	for i := len(r.resultStack) - 1; i >= 0; i-- {
		if err := r.resultStack[i](ctx); err != nil {
			return fmt.Errorf("rollback step with index %d step: %w", i, err)
		}
	}

	return nil
}

func (r *Runner) runStep(ctx context.Context, step Step[any]) (any, error) {
	res, err := retry[any](func() (any, error) {
		return step.Run(ctx)
	}, step.ShouldRetryRun)
	if err != nil {
		return nil, fmt.Errorf("run step via runner: %w", err)
	}

	// In case, step has rollback action, we should save it in runner callstack
	// for future usage.
	if step.Rollback != nil {
		r.resultStack = append(r.resultStack, func(ctx context.Context) error {
			return retryWithoutResult(func() error {
				return step.Rollback(ctx, res)
			}, step.ShouldRetryRollback)
		})
	}

	return res, err
}
