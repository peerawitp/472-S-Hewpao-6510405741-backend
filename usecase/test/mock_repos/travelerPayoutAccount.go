package mock_repos

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockTravelerPayoutAccountRepository struct {
	mock.Mock
}

func (m *MockTravelerPayoutAccountRepository) Store(ctx context.Context, travelerPayoutAccount *domain.TravelerPayoutAccount) error {
	args := m.Called(ctx, travelerPayoutAccount)
	return args.Error(0)
}

func (m *MockTravelerPayoutAccountRepository) FindByUserID(ctx context.Context, userID string) ([]domain.TravelerPayoutAccount, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.TravelerPayoutAccount), args.Error(1)
}
