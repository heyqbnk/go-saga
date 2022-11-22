package main

import (
	"fmt"
	"github.com/heyqbnk/transaction"
	"github.com/pkg/errors"
)

type TransferTransaction struct {
	trx              *transaction.Transaction
	From, To, Amount int64
}

func (t *TransferTransaction) Perform() error {
	return t.trx.Perform()
}

func (t *TransferTransaction) PreRPC() error {
	return nil
}

func (t *TransferTransaction) RPC(runner transaction.StepRunner) {
	trxID := transaction.RunStep(runner, &transaction.Step[int]{
		Run: func() (int, error) {
			return 920, nil
		},
		RunRetries: 3,
		Rollback: func(trxID int) error {
			fmt.Println("Attempt to rollback", trxID)
			return nil
		},
		RollbackRetries: 3,
	})

	_ = transaction.RunStep(runner, &transaction.Step[any]{
		Run: func() (any, error) {
			fmt.Println("I know about trx id", trxID)
			return nil, errors.New("SOME ERROR")
		},
	})
}

func (t *TransferTransaction) PostRPC() error {
	return nil
}

func (t *TransferTransaction) MaxRetries() int {
	return 0
}

func (t *TransferTransaction) Lock() error {
	return nil
}

func (t *TransferTransaction) Unlock() error {
	return nil
}

func NewTransferTransaction(from, to, amount int64) *TransferTransaction {
	t := &TransferTransaction{
		From:   from,
		To:     to,
		Amount: amount,
	}
	t.trx = transaction.New(t)
	return t
}

func main() {
	trx := NewTransferTransaction(3, 5, 10_000)
	fmt.Println(trx.Perform())
}
