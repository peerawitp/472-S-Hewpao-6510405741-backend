package mock_repos

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Store(ctx context.Context, transaction *domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByThirdPartyPaymentID(ctx context.Context, thirdPartyPaymentID string) (*domain.Transaction, error) {
	args := m.Called(ctx, thirdPartyPaymentID)
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) UpdateStatusByThirdPartyPaymentID(ctx context.Context, thirdPartyPaymentID string, status types.PaymentStatus) error {
	args := m.Called(ctx, thirdPartyPaymentID, status)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}
