package buyitemsaga

import (
	"context"

	"github.com/heyqbnk/go-saga/example/dto"
)

// Wallet is a service which manipulates user balance.
type Wallet interface {
	DecreaseUserBalance(ctx context.Context, userID, amount int) error
	IncreaseUserBalance(ctx context.Context, userID, amount int) error
}

// UserItems is a dependency which manipulates user inventory.
type UserItems interface {
	GiveItem(ctx context.Context, userID, itemID int) (recordID int, err error)
	RevokeItem(ctx context.Context, recordID int) error
}

// Store is a dependency which manipulates shop items.
type Store interface {
	GetItemByID(ctx context.Context, itemID int) (dto.ShopItem, error)
	IncreasePurchaseCounter(ctx context.Context, itemID int) error
	DecreasePurchaseCounter(ctx context.Context, itemID int) error
}

// Locker is a dependency which allows locking sagas.
type Locker interface {
	Lock(ctx context.Context, key string) error
	Unlock(ctx context.Context, key string) error
}
