package entity

import (
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	ID        string
	AccountID string
	Balance   float64
	CreatedAt time.Time
}

func NewBalance(accountId string, balance float64) *Balance {
	return &Balance{
		ID:        uuid.New().String(),
		AccountID: accountId,
		Balance:   balance,
		CreatedAt: time.Now(),
	}
}
