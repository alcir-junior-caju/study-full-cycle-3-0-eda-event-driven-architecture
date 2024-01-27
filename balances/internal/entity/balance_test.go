package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBalance(t *testing.T) {
	balance := NewBalance("123", 100.00)
	assert.NotNil(t, balance)
	assert.Equal(t, "123", balance.AccountID)
	assert.Equal(t, 100.00, balance.Balance)
}
