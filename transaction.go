package gotransaction

import (
	"context"
)

type Transaction interface {
	Run(ctx context.Context, runner *Runner) error
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	ShouldRetryUnlock(err error, callsCount int) bool
}

// RunTransaction runs whole transaction. Steps:
// 1. Lock transaction.
// 2. Run all steps.
// 3. Unlock transaction.
func RunTransaction(ctx context.Context, trx Transaction) (err error) {
	if err = trx.Lock(ctx); err != nil {
		return err
	}

	defer func() {
		unlockErr := retryWithoutResult(func() error {
			return trx.Unlock(ctx)
		}, trx.ShouldRetryUnlock)

		if unlockErr != nil {
			if err == nil {
				err = unlockErr
			} else {
				// FIXME: What to do with such type of error.
			}
		}
	}()

	if err = trx.Run(ctx, newRunner()); err != nil {
		return err
	}

	return nil
}
