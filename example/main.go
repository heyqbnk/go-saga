package main

import (
	"context"
	"fmt"
	"log"
	"time"

	buyitemsaga "github.com/heyqbnk/go-saga/example/buy-item-saga"
	examplelocker "github.com/heyqbnk/go-saga/example/locker"
	examplestore "github.com/heyqbnk/go-saga/example/store"
	exampleuseritems "github.com/heyqbnk/go-saga/example/user-items"
	examplewallet "github.com/heyqbnk/go-saga/example/wallet"
)

func main() {
	lock := examplelocker.New()
	store := examplestore.New()
	userItems := exampleuseritems.New()
	wallet := examplewallet.New()

	saga := buyitemsaga.New(lock, store, userItems, wallet)

	execCtx, _ := context.WithTimeout(context.Background(), time.Second)
	rollbackCtx := context.Background()

	res, execErr, rollbackErr, unlockErr := saga.Run(execCtx, rollbackCtx, 777, 90)
	if execErr != nil {
		log.Fatal(fmt.Errorf("execution: %w", execErr))
	}
	if rollbackErr != nil {
		log.Fatal(fmt.Errorf("rollback: %w", rollbackErr))
	}
	if unlockErr != nil {
		log.Fatal(fmt.Errorf("unlock: %w", unlockErr))
	}

	log.Printf("saga result: %s", res)
}
