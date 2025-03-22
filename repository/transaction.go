package repository

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/types"
)

type TransactionRepository interface {
	Store(ctx context.Context, transaction *domain.Transaction) error
	FindByID(ctx context.Context, id string) (*domain.Transaction, error)
	FindByThirdPartyPaymentID(ctx context.Context, thirdPartyPaymentID string) (*domain.Transaction, error)
	UpdateStatusByThirdPartyPaymentID(ctx context.Context, thirdPartyPaymentID string, status types.PaymentStatus) error
	FindByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error)
}
