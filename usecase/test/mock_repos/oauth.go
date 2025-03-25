package mock_repos

import (
	"context"

	"github.com/hewpao/hewpao-backend/dto"
	"github.com/stretchr/testify/mock"
)

type MockOAuthRepository struct {
	mock.Mock
}

func (m *MockOAuthRepository) VerifyIDToken(ctx context.Context, idToken string) (*dto.OAuthClaims, error) {
	args := m.Called(ctx, idToken)
	return args.Get(0).(*dto.OAuthClaims), args.Error(1)
}
