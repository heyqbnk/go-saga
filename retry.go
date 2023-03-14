package gotransaction

// Retries "f" until "shouldRetry" returns true. In case, "shouldRetry" is
// nil, retry calls "f" function only once.
func retry[T interface{}](
	f func() (T, error),
	shouldRetry func(err error, callsCount int) bool,
) (res T, err error) {
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
func retryWithoutResult(f func() error, shouldRetry func(err error, callsCount int) bool) error {
	_, err := retry[any](func() (any, error) {
		return nil, f()
	}, shouldRetry)

	return err
}
