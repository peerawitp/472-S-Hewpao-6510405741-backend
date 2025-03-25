package mock_repos

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockBankRepository struct {
	mock.Mock
}

func (m *MockBankRepository) GetBySwift(swift string) (*domain.Bank, error) {
	args := m.Called(swift)
	return args.Get(0).(*domain.Bank), args.Error(1)
}

func (m *MockBankRepository) GetAll() ([]domain.Bank, error) {
	args := m.Called()
	return args.Get(0).([]domain.Bank), args.Error(1)
}
