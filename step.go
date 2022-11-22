package transaction

type Step[T interface{}] struct {
	// Run runs this step.
	Run func() (T, error)
	// RunRetries contains count of retries for Run step.
	RunRetries int

	// Rollback rollbacks step operation.
	Rollback func(value T) error
	// RollbackRetries contains count of retries for Rollback step.
	RollbackRetries int
}

// StepRunner represents function which accepts step description and its
// result destination pointer.
type StepRunner = func(step *Step[any], dest *any)

// RunStep accepts step information and saves it in passed store which is
// usually provided by transaction.
func RunStep[T interface{}](runner StepRunner, step *Step[T]) T {
	var result interface{}
	runner(&Step[any]{
		Run: func() (any, error) {
			return step.Run()
		},
		RunRetries: step.RunRetries,
		Rollback: func(value any) error {
			v, _ := value.(T)
			return step.Rollback(v)
		},
		RollbackRetries: step.RollbackRetries,
	}, &result)

	v, _ := result.(T)
	return v
}
