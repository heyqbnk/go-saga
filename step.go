package gotransaction

import "context"

type Step[T interface{}] struct {
	// Run is a function which executes current step.
	Run func(ctx context.Context) (T, error)

	// Rollback is a function which rollbacks current step.
	Rollback func(ctx context.Context, runResult T) error

	// ShouldRetryRun is a function which should return true in case,
	// Run step should be retried.
	ShouldRetryRun func(err error, callsCount int) bool

	// ShouldRetryRollback is a function which should return true in case,
	// Rollback step should be retried.
	ShouldRetryRollback func(err error, callsCount int) bool
}

func RunStep[T interface{}](ctx context.Context, runner *Runner, step Step[T]) (result T, err error) {
	if step.Run == nil {
		return result, ErrStepRunMissing
	}

	newStep := Step[interface{}]{
		Run: func(ctx context.Context) (interface{}, error) {
			return step.Run(ctx)
		},
		ShouldRetryRun:      step.ShouldRetryRun,
		ShouldRetryRollback: step.ShouldRetryRollback,
	}

	if step.Rollback != nil {
		newStep.Rollback = func(ctx context.Context, runResult interface{}) error {
			stepResult, ok := runResult.(T)
			if !ok {
				return ErrIncorrectResultType
			}

			return step.Rollback(ctx, stepResult)
		}
	}

	value, err := runner.runStep(ctx, newStep)
	if err != nil {
		return result, err
	}

	if value == nil {
		return result, nil
	}

	result, ok := value.(T)
	if !ok {
		return result, ErrIncorrectResultType
	}

	return result, nil
}
