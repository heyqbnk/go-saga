package buyitemsaga

import (
	"context"

	gosaga "github.com/heyqbnk/go-saga"
)

type Factory struct {
	locker    Locker
	store     Store
	userItems UserItems
	wallet    Wallet
}

func New(locker Locker, store Store, userItems UserItems, wallet Wallet) *Factory {
	return &Factory{locker: locker, store: store, userItems: userItems, wallet: wallet}
}

func (s *Factory) Run(
	execCtx, rollbackCtx context.Context,
	userID, itemID int,
) (res string, execErr, rollbackErr, unlockErr error) {
	return gosaga.RunSaga[string](execCtx, rollbackCtx, &saga{
		itemID:    itemID,
		locker:    s.locker,
		store:     s.store,
		userItems: s.userItems,
		userID:    userID,
		wallet:    s.wallet,
	})
}
