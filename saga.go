package saga

import (
	"context"
	"fmt"
)

type Saga[T any] interface {
	// Run runs current saga.
	Run(execCtx, rollbackCtx context.Context, runner *Runner) (result T, execErr, rollbackErr error)
	// Lock locks this saga and prevents its parallel execution.
	Lock(ctx context.Context) error
	// ShouldRetryUnlock is the function which determines if saga unlock should
	// be retried again. See ShouldRetryFunc.
	ShouldRetryUnlock(err error, retriesCount int) bool
	// Unlock unlocks saga.
	Unlock(ctx context.Context) error
}

// RunSaga runs specified saga.
func RunSaga[T any](
	execCtx, rollbackCtx context.Context,
	saga Saga[T],
) (result T, execErr, rollbackErr, unlockErr error) {
	// Lock saga.
	if execErr = saga.Lock(execCtx); execErr != nil {
		return result, fmt.Errorf("lock saga: %w", execErr), nil, nil
	}

	// Don't forget to unlock saga.
	defer func() {
		unlockErr = retryWithoutResult(func() error {
			return saga.Unlock(execCtx)
		}, saga.ShouldRetryUnlock)
	}()

	// Execute saga.
	result, execErr, rollbackErr = saga.Run(execCtx, rollbackCtx, newRunner())
	if execErr != nil {
		execErr = fmt.Errorf("run saga: %w", execErr)
	}
	if rollbackErr != nil {
		rollbackErr = fmt.Errorf("rollback saga: %w", rollbackErr)
	}

	return result, execErr, rollbackErr, nil
}
