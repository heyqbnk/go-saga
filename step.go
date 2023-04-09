package saga

import "context"

// Step describes step in saga.
type Step[T any] struct {
	// Run is a function which executes current step.
	Run func(ctx context.Context) (T, error)

	// Rollback is a function which rollbacks current step.
	Rollback func(ctx context.Context, runResult T) error

	// ShouldRetryRun is a function which should return true in case,
	// Run step should be retried.
	//
	// This method is optional. In case, function was not specified, runner
	// will not try to retry Run method.
	ShouldRetryRun ShouldRetryFunc

	// ShouldRetryRollback is a function which should return true in case,
	// Rollback step should be retried.
	//
	// This method is optional. In case, function was not specified, runner
	// will not try to retry Rollback method.
	ShouldRetryRollback ShouldRetryFunc
}

// RunStep executes step. Panic will be called in case, step's Run method
// is not specified.
func RunStep[T any](ctx context.Context, runner *Runner, step Step[T]) (result T, err error) {
	if step.Run == nil {
		panic("step is missing Run method")
	}

	// Create step which return unknown result.
	newStep := Step[any]{
		Run: func(ctx context.Context) (any, error) {
			return step.Run(ctx)
		},
		ShouldRetryRun:      step.ShouldRetryRun,
		ShouldRetryRollback: step.ShouldRetryRollback,
	}

	// In case, rollback action is specified, we should copy it.
	if step.Rollback != nil {
		newStep.Rollback = func(ctx context.Context, runResult any) error {
			stepResult, _ := runResult.(T)
			return step.Rollback(ctx, stepResult)
		}
	}

	// Run step via runner.
	value, err := runner.runStep(ctx, newStep)
	if err != nil {
		return result, err
	}

	// Special case, value has the pointer type.
	if value == nil {
		return result, nil
	}

	// Convert result to T type.
	result, _ = value.(T)

	return result, nil
}
