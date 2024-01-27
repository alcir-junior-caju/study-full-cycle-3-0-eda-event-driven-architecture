package find_account

import (
	"testing"

	"github.com.br/devfullcycle/fc-ms-balances/internal/entity"
	"github.com.br/devfullcycle/fc-ms-balances/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAccountUseCase_Execute(t *testing.T) {
	balance := entity.NewBalance("1", 10.0)
	accountMock := &mocks.BalanceGatewayMock{}
	accountMock.On("FindByAccountID", balance.ID).Return(balance, nil)

	balanceMock := &mocks.BalanceGatewayMock{}
	balanceMock.On("Save", mock.Anything).Return(nil)

	uc := NewFindAccountUseCase(accountMock)
	inputDto := FindAccountInputDTO{
		AccountID: balance.ID,
	}
	output, err := uc.Execute(inputDto)
	assert.Nil(t, err)
	assert.Equal(t, balance.ID, output.ID)
	assert.Equal(t, balance.AccountID, output.AccountID)
	assert.Equal(t, balance.Balance, output.Balance)
	assert.Equal(t, balance.CreatedAt, output.CreatedAt)
}
