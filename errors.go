package gotransaction

import "errors"

var (
	ErrStepRunMissing      = errors.New("run method of step missing")
	ErrIncorrectResultType = errors.New("incorrect result type")
)
