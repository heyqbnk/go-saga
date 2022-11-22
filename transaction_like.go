package transaction

// TransactionLike describes minimal set of methods, required by Transaction
// structure.
type TransactionLike interface {
	// PreRPC represents a step which is called before main step.
	//
	// You are unable to roll back this step, because it is assumed, that you
	// use this step to check requirements and other. We call it read-only step.
	PreRPC() error

	// RPC represents main transaction step. This step can contain any data
	// mutations as long as in case of error, they could be rolled back.
	//
	// Function accepts step runner which is able to run passed steps.
	RPC(runner StepRunner)

	// PostRPC is step which is called after main transaction step is completed.
	// This step should also not contain any mutations or critical actions.
	PostRPC() error

	// MaxRetries returns maximum count of retries for PreRPC and PostRPC steps.
	// RPC step has its own sub-steps with their own retries description.
	MaxRetries() int

	// Lock locks transaction and does not allow to call it several times.
	Lock() error

	// Unlock unlocks transaction execution.
	Unlock() error
}
