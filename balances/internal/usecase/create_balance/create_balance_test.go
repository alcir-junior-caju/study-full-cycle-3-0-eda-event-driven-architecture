package create_balance

import (
	"testing"

	"github.com.br/devfullcycle/fc-ms-balances/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBalanceUseCase_Execute(t *testing.T) {
	balanceMock := &mocks.BalanceGatewayMock{}
	balanceMock.On("Save", mock.Anything).Return(nil)

	uc := NewCreateBalanceUseCase(balanceMock)
	inputDto := CreateBalanceInputDTO{
		AccountID: "1",
		Balance:   100.00,
	}
	output, err := uc.Execute(inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	// asssert valid uuid
	balanceMock.AssertExpectations(t)
	balanceMock.AssertNumberOfCalls(t, "Save", 1)
}
