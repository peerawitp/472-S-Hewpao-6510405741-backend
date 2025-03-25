package mock_repos

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) CreatePayment(ctx context.Context, pr *domain.ProductRequest) (*dto.CreatePaymentResponseDTO, error) {
	args := m.Called(ctx, pr)
	return args.Get(0).(*dto.CreatePaymentResponseDTO), args.Error(1)
}

type MockPaymentRepositoryFactory struct {
	mock.Mock
}

func (m *MockPaymentRepositoryFactory) Register(provider string, repo repository.PaymentRepository) {
	m.Called(provider, repo)
}

func (m *MockPaymentRepositoryFactory) GetRepository(provider string) (repository.PaymentRepository, error) {
	args := m.Called(provider)
	if repo, ok := args.Get(0).(repository.PaymentRepository); ok {
		return repo, args.Error(1)
	}
	return nil, args.Error(1)
}
