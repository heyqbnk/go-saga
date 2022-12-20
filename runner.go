package transaction

import (
	"github.com/pkg/errors"
)

type cachedStep[T interface{}] struct {
	step   Step[T]
	result T
}

// Runner is passed to transaction to work with transaction state.
type Runner struct {
	// List of completed steps.
	cachedSteps []cachedStep[any]
}

func newRunner() *Runner {
	return &Runner{}
}

// Run runs specified step and caches it for future rollback.
func (r *Runner) Run(step Step[any]) (res any, err error) {
	if step.Run == nil {
		return nil, errors.New("\"Run\" function should be specified")
	}
	err = retry(func() error {
		res, err = step.Run()
		return err
	}, step.Retries+1)
	if err != nil {
		return nil, err
	}

	// Cache step result.
	r.cacheStep(cachedStep[any]{step: step, result: res})

	return res, err
}

// Rollback rollbacks transaction.
func (r *Runner) Rollback() error {
	for i := len(r.cachedSteps) - 1; i >= 0; i-- {
		rollback := r.cachedSteps[i].step.Rollback
		if rollback == nil {
			continue
		}
		if err := rollback(r.cachedSteps[i].result); err != nil {
			return err
		}
	}
	return nil
}

// Caches step.
func (r *Runner) cacheStep(step cachedStep[any]) {
	r.cachedSteps = append(r.cachedSteps, step)
}
