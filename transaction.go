package transaction

import "context"

type Transaction interface {
	Run(ctx context.Context, runner *Runner) (res any, err error)

	// Lock locks this transaction and does not allow to execute it several
	// times.
	Lock(ctx context.Context) error

	// Unlock unlocks transaction for execution.
	Unlock(ctx context.Context) error

	// MaxRetries returns count of retries applied to these methods:
	// - Unlock
	MaxRetries() int
}
