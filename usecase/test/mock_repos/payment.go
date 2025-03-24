package mock_repos

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) CreatePayment(ctx context.Context, pr *domain.ProductRequest) (*dto.CreatePaymentResponseDTO, error) {
	args := m.Called(ctx, pr)
	return args.Get(0).(*dto.CreatePaymentResponseDTO), args.Error(1)
}
