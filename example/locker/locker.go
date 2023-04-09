package locker

import (
	"context"
	"log"
)

type Locker struct {
}

func New() *Locker {
	return &Locker{}
}

func (l *Locker) Lock(ctx context.Context, key string) error {
	log.Printf("locked %s key", key)
	return nil
}

func (l *Locker) Unlock(ctx context.Context, key string) error {
	log.Printf("unlocked %s key", key)
	return nil
}
