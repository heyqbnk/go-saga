package transaction

// Calls passed function count times or until it returned no error.
func retry(f func() error, count int) (err error) {
	for i := 0; i < count; i++ {
		if err = f(); err == nil {
			break
		}
	}
	return err
}
