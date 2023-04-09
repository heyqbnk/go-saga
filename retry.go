package saga

// ShouldRetryFunc function which determines if some function should be called
// again. Accepts last occurred error and already performed retries count.
type ShouldRetryFunc = func(err error, retriesCount int) bool

// ShouldRetriesCount return function which retries function execution not more
// than count **additional** times.
//
// For example, when you specify 3 here, this will lead to maximum 4 function
// calls - 1 guaranteed + 3 additional.
func ShouldRetriesCount(count int) ShouldRetryFunc {
	return func(_ error, callsCount int) bool {
		return callsCount < count
	}
}

// Retries specified function until retrier function is returning true.
// In case, retrier is nil, specified function will only be called once.
func retry[T any](f func() (T, error), shouldRetry ShouldRetryFunc) (res T, err error) {
	retriesCount := 0

	for {
		// Execute wrapped function.
		res, err = f()
		if err == nil {
			return res, nil
		}

		if shouldRetry == nil || !shouldRetry(err, retriesCount) {
			return res, err
		}

		// Otherwise, increase retries count and try again.
		retriesCount++
	}
}

// See retry.
func retryWithoutResult(f func() error, shouldRetry ShouldRetryFunc) error {
	_, err := retry[any](func() (any, error) {
		return nil, f()
	}, shouldRetry)

	return err
}
