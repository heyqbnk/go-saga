package useritems

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type UserItems struct {
	random *rand.Rand
}

func New() *UserItems {
	return &UserItems{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (u *UserItems) GiveItem(ctx context.Context, userID, itemID int) (recordID int, err error) {
	log.Printf("giving user %d item %d", userID, itemID)
	return u.random.Intn(100000), nil
}

func (u *UserItems) RevokeItem(ctx context.Context, recordID int) error {
	log.Printf("revoking user item by record ID %d", recordID)
	return nil
}
