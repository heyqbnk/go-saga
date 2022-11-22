package transaction

import (
	"github.com/pkg/errors"
)

var (
	ErrPreRPCFailed   = errors.New("preRPC step failed")
	ErrPostRPCFailed  = errors.New("postRPC step failed")
	ErrUnableToLock   = errors.New("unable to lock")
	ErrUnableToUnlock = errors.New("unable to unlock")
)

type Transaction struct {
	TransactionLike
}

// Perform performs transaction with all its steps.
func (t *Transaction) Perform() error {
	// Perform preRPC step.
	if err := t.preRPC(); err != nil {
		if errors.Is(err, ErrUnableToLock) {
			return err
		}

		// Something went wrong, roll this step back.
		if err := t.rollbackPreRPC(); err != nil {
			// TODO: What should we do here?
		}

		// TODO: What if we were unable to rollback? We should log at least.
		return errors.Wrap(ErrPreRPCFailed, err.Error())
	}

	// Perform RPC step.
	if err := t.rpc(); err != nil {
		return err
	}

	// Perform post RPC step.
	if err := t.postRPC(); err != nil {
		return errors.Wrap(ErrPostRPCFailed, err.Error())
	}
	return nil
}

// Performs pre RPC step of transaction. It locks, transaction, and then,
// executes passed PreRPC step.
func (t *Transaction) preRPC() error {
	// Try to lock transaction.
	if err := t.Lock(); err != nil {
		return errors.Wrap(ErrUnableToLock, err.Error())
	}
	return retry(t.PreRPC, t.MaxRetries())
}

// Rollbacks pre RPC step unlocking transaction.
func (t *Transaction) rollbackPreRPC() error {
	if err := retry(t.Unlock, t.MaxRetries()); err != nil {
		return errors.Wrap(ErrUnableToUnlock, err.Error())
	}
	return nil
}

// Performs RPC step. It runs all transaction steps with passed step runner.
// Function returns error, returned by one of defined by RPC() steps.
func (t *Transaction) rpc() error {
	// Create store, which caches steps results.
	store := newStepStore()

	// Last saved step error. We need this value to return from this function
	// to let executor know which error occurred.
	var err error

	// Create RPC step runner. It will be passed to transaction RPC() method.
	// Internally, this function should accept information about step and
	// cache its result in case, step was completed. In case, transaction
	// failed, runner should roll back whole transaction by specified Rollback()
	// functions.
	runner := func(step *Step[any], dest *any) {
		var res any

		// Run "Run" step
		for i := 0; i < step.RunRetries+1; i++ {
			res, err = step.Run()
			if err == nil {
				break
			}
		}

		// In case, we were unable to proceed this step, we should roll back all
		// previous steps.
		if err != nil {
			// Get last stored step result and roll it back.
			for {
				cached := store.Pop()
				if cached == nil {
					break
				}
				for i := 0; i < cached.step.RollbackRetries+1; i++ {
					// While rolling back, we always should pass step result.
					if rbErr := cached.step.Rollback(cached.res); rbErr != nil {
						err = rbErr
					}
				}
			}
			return
		}

		// Otherwise, we memoize step result for future steps.
		store.Store(step, res)

		// And store result in specified destination.
		*dest = res
	}

	// Run RPC step with runner.
	t.RPC(runner)

	return err
}

// Performs post RPC step.
func (t *Transaction) postRPC() error {
	retries := t.MaxRetries()

	// Perform transaction's post rpc step.
	err := retry(t.PostRPC, retries)

	// Unlock transaction.
	if lockErr := retry(t.Unlock, retries); lockErr != nil {
		// TODO: What should we do?
	}
	return err
}

func New(transaction TransactionLike) *Transaction {
	return &Transaction{transaction}
}
