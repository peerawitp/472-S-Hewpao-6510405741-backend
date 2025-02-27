package repository

import (
	"context"
	"github.com/hewpao/hewpao-backend/domain"
)

type TransactionRepository interface {
	Store(ctx context.Context, transaction *domain.Transaction) error
	FindByID(ctx context.Context, id string) (*domain.Transaction, error)
}





