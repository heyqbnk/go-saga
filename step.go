package transaction

// Step represents transaction step.
type Step[T interface{}] struct {
	// Run runs this step.
	Run func() (res T, err error)
	// Rollback rollbacks this step accepting Run result.
	Rollback func(res T) error
	// Retries contains count of retries for Run and Rollback.
	Retries int
}
