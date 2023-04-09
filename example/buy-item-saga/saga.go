package buyitemsaga

import (
	"context"
	"fmt"

	gosaga "github.com/heyqbnk/go-saga"
	"github.com/heyqbnk/go-saga/example/dto"
)

var (
	// Default retrier function.
	_shouldRetry = gosaga.ShouldRetriesCount(3)
)

// Implements gosaga.Saga interface.
type saga struct {
	itemID    int
	locker    Locker
	store     Store
	userItems UserItems
	userID    int
	wallet    Wallet
}

func (s *saga) Run(
	execCtx, rollbackCtx context.Context,
	runner *gosaga.Runner,
) (res string, execErr, rollbackErr error) {
	shopItem, err := s.getShopItem(execCtx)
	if err != nil {
		return "", fmt.Errorf("get shop item: %w", err), nil
	}

	// 1. Decrease user balance
	if err := s.stepDecreaseBalance(execCtx, runner, shopItem.Price); err != nil {
		return "", fmt.Errorf("step decrease user balance: %w", err), nil
	}

	// Don't forget to rollback saga in case, something went wrong.
	defer func() {
		if execErr != nil {
			rollbackErr = runner.Rollback(rollbackCtx)
		}
	}()

	// 2. Give user an item.
	if err := s.stepGiveItem(execCtx, runner); err != nil {
		return "", fmt.Errorf("step collect user item: %w", err), nil
	}

	// 3. Increase purchase counter.
	if err := s.stepIncreasePurchaseCounter(execCtx, runner); err != nil {
		return "", fmt.Errorf("step increase purchase counter: %w", err), nil
	}

	return "completed successfully", nil, nil
}

func (s *saga) Lock(ctx context.Context) error {
	if err := s.locker.Lock(ctx, s.lockKey()); err != nil {
		return fmt.Errorf("lock transaction: %w", err)
	}

	return nil
}

func (s *saga) ShouldRetryUnlock(_ error, callsCount int) bool {
	// We don't want saga unlock to retry more than 3 times.
	return callsCount < 3
}

func (s *saga) Unlock(ctx context.Context) error {
	if err := s.locker.Unlock(ctx, s.lockKey()); err != nil {
		return fmt.Errorf("lock transaction: %w", err)
	}

	return nil
}

// Returns related shop item information.
func (s *saga) getShopItem(ctx context.Context) (dto.ShopItem, error) {
	item, err := s.store.GetItemByID(ctx, s.itemID)
	if err != nil {
		return dto.ShopItem{}, fmt.Errorf("get by id: %w", err)
	}

	return item, nil
}

// Returns locker key.
func (s *saga) lockKey() string {
	return fmt.Sprintf("buy_item_%d_%d", s.userID, s.itemID)
}

// Step which decreases user balance.
func (s *saga) stepDecreaseBalance(ctx context.Context, runner *gosaga.Runner, amount int) error {
	_, err := gosaga.RunStep(ctx, runner, gosaga.Step[any]{
		Run: func(ctx context.Context) (any, error) {
			if err := s.wallet.DecreaseUserBalance(ctx, s.userID, amount); err != nil {
				return nil, fmt.Errorf("decrease balance via wallet: %w", err)
			}

			return nil, nil
		},
		Rollback: func(ctx context.Context, _ any) error {
			if err := s.wallet.IncreaseUserBalance(ctx, s.userID, amount); err != nil {
				return fmt.Errorf("increase balance via wallet: %w", err)
			}

			return nil
		},
		ShouldRetryRun:      _shouldRetry,
		ShouldRetryRollback: _shouldRetry,
	})
	if err != nil {
		return fmt.Errorf("run step: %w", err)
	}

	return nil
}

// Step which gives user an item.
func (s *saga) stepGiveItem(ctx context.Context, runner *gosaga.Runner) error {
	_, err := gosaga.RunStep(ctx, runner, gosaga.Step[int]{
		Run: func(ctx context.Context) (int, error) {
			recordID, err := s.userItems.GiveItem(ctx, s.userID, s.itemID)
			if err != nil {
				return 0, fmt.Errorf("give item via user items: %w", err)
			}

			return recordID, nil
		},
		Rollback: func(ctx context.Context, recordID int) error {
			if err := s.userItems.RevokeItem(ctx, recordID); err != nil {
				return fmt.Errorf("revoke item via user items: %w", err)
			}

			return nil
		},
		ShouldRetryRun:      _shouldRetry,
		ShouldRetryRollback: _shouldRetry,
	})
	if err != nil {
		return fmt.Errorf("run step: %w", err)
	}

	return nil
}

// Step which increases item purchase counter.
func (s *saga) stepIncreasePurchaseCounter(ctx context.Context, runner *gosaga.Runner) error {
	_, err := gosaga.RunStep(ctx, runner, gosaga.Step[any]{
		Run: func(ctx context.Context) (any, error) {
			if err := s.store.IncreasePurchaseCounter(ctx, s.itemID); err != nil {
				return nil, fmt.Errorf("increase purchase counter via store: %w", err)
			}

			return nil, nil
		},
		Rollback: func(ctx context.Context, _ any) error {
			if err := s.store.DecreasePurchaseCounter(ctx, s.itemID); err != nil {
				return fmt.Errorf("decrease purchase counter via store: %w", err)
			}

			return nil
		},
		ShouldRetryRun:      _shouldRetry,
		ShouldRetryRollback: _shouldRetry,
	})
	if err != nil {
		return fmt.Errorf("run step: %w", err)
	}

	return nil
}
