package wallet

import (
	"context"
	"log"
)

type Wallet struct {
}

func New() *Wallet {
	return &Wallet{}
}

func (w *Wallet) DecreaseUserBalance(ctx context.Context, userID, amount int) error {
	log.Printf("decreasing user %d balance by %d", userID, amount)
	return nil
}

func (w *Wallet) IncreaseUserBalance(ctx context.Context, userID, amount int) error {
	log.Printf("increasing user %d balance by %d", userID, amount)
	return nil
}
