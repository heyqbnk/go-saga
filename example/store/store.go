package store

import (
	"context"
	"log"

	"github.com/heyqbnk/go-saga/example/dto"
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

func (s *Store) DecreasePurchaseCounter(ctx context.Context, itemID int) error {
	log.Printf("decreasing purchase counter of item with ID %d", itemID)
	return nil
}

func (s *Store) GetItemByID(ctx context.Context, itemID int) (dto.ShopItem, error) {
	log.Printf("getting shop item by ID %d", itemID)

	return dto.ShopItem{
		Price: 15499,
	}, nil
}

func (s *Store) IncreasePurchaseCounter(ctx context.Context, itemID int) error {
	log.Printf("increasing purchase counter of item with ID %d", itemID)
	return nil
}
