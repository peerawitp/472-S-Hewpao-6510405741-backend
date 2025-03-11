package repository

import (
	"context"
	"errors"
	"log"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, pr *domain.ProductRequest) (*dto.CreatePaymentResponseDTO, error)
}

type PaymentRepositoryFactory struct {
	repos map[string]PaymentRepository
}

func NewPaymentRepositoryFactory() PaymentRepositoryFactory {
	return PaymentRepositoryFactory{repos: make(map[string]PaymentRepository)}
}

func (f *PaymentRepositoryFactory) Register(provider string, repo PaymentRepository) {
	log.Println("[💸] Registered", provider, "Payment repository!")
	f.repos[provider] = repo
}

func (f *PaymentRepositoryFactory) GetRepository(provider string) (PaymentRepository, error) {
	repo, exists := f.repos[provider]
	if !exists {
		return nil, errors.New("unsupported Payment provider")
	}
	return repo, nil
}
