package transaction

import "context"

// Run runs transaction and returns parameters, returned by t.Run method.
func Run(ctx context.Context, t Transaction) (res any, err error) {
	// Lock transaction.
	if err = t.Lock(ctx); err != nil {
		return nil, err
	}

	// Don't forget to unlock transaction.
	defer func() {
		unlockErr := retry(func() error {
			return t.Unlock(ctx)
		}, t.MaxRetries())
		if unlockErr != nil {
			err = unlockErr
		}
	}()

	return t.Run(ctx, newRunner())
}

// RunStepWithRunner runs step with specified runner and returns type T. In
// other words, this function is just better typed runner.Run function.
func RunStepWithRunner[T interface{}](
	runner *Runner,
	step Step[T],
) (res T, err error) {
	resInf, err := runner.Run(Step[any]{
		Run: func() (any, error) {
			return step.Run()
		},
		Rollback: func(res any) error {
			t, _ := res.(T)
			return step.Rollback(t)
		},
		Retries: 0,
	})
	res, _ = resInf.(T)
	return res, err
}

// Calls passed function count times or until it returned no error.
func retry(f func() error, count int) (err error) {
	for i := 0; i < count; i++ {
		if err = f(); err == nil {
			break
		}
	}
	return err
}
