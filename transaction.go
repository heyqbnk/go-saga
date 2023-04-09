package gotransaction

import (
	"context"
	"fmt"
)

type Transaction interface {
	Run(execCtx, rollbackCtx context.Context, runner *Runner) (execErr, rollbackErr error)
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	ShouldRetryUnlock(err error, callsCount int) bool
}

// RunTransaction runs whole transaction. Steps:
// 1. Lock transaction.
// 2. Run all steps.
// 3. Unlock transaction.
func RunTransaction(
	execCtx, rollbackCtx context.Context,
	trx Transaction,
) (execErr, rollbackErr, unlockErr error) {
	if execErr = trx.Lock(execCtx); execErr != nil {
		return fmt.Errorf("lock transaction: %w", execErr), nil, nil
	}

	defer func() {
		unlockErr = retryWithoutResult(func() error {
			return trx.Unlock(execCtx)
		}, trx.ShouldRetryUnlock)
	}()

	execErr, rollbackErr = trx.Run(execCtx, rollbackCtx, newRunner())
	if execErr != nil {
		execErr = fmt.Errorf("run transaction: %w", execErr)
	}
	if rollbackErr != nil {
		rollbackErr = fmt.Errorf("rollback transaction: %w", rollbackErr)
	}

	return execErr, rollbackErr, nil
}
