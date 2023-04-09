package gotransaction

// ShouldRetryFunc function which determines if some function should be called
// again. Accepts last occurred error and current function calls count.
type ShouldRetryFunc = func(err error, callsCount int) bool

// ShouldRetriesCount return function which retries function execution not more
// than count times.
func ShouldRetriesCount(count int) ShouldRetryFunc {
	return func(_ error, callsCount int) bool {
		return callsCount <= count
	}
}

// Retries "f" until "shouldRetry" returns true. In case, "shouldRetry" is
// nil, retry calls "f" function only once.
func retry[T interface{}](f func() (T, error), shouldRetry ShouldRetryFunc) (res T, err error) {
	doRetry := true
	callsCount := 0

	for doRetry {
		callsCount++
		res, err = f()
		if err != nil {
			if shouldRetry == nil {
				doRetry = false
			} else {
				doRetry = shouldRetry(err, callsCount)
			}
			continue
		}
		doRetry = false
	}

	return res, err
}

// See retry.
func retryWithoutResult(f func() error, shouldRetry ShouldRetryFunc) error {
	_, err := retry[any](func() (any, error) {
		return nil, f()
	}, shouldRetry)

	return err
}
