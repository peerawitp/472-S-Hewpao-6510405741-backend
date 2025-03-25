package mock_repos

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockVerificationRepository struct {
	mock.Mock
}

func (m *MockVerificationRepository) Create(verification *domain.Verification) error {
	args := m.Called(verification)
	return args.Error(0)
}

func (m *MockVerificationRepository) FindByID(id uint) (*domain.Verification, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Verification), args.Error(1)
}
