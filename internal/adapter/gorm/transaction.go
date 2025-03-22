package gorm

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/types"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Store(ctx context.Context, transaction *domain.Transaction) error {
	if err := r.db.Where("product_request_id = ? AND status = ?", transaction.ProductRequestID, types.PaymentSuccess).First(&domain.Transaction{}).Error; err == nil {
		return exception.ErrProductRequestAlreadyPaid
	}

	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) FindByThirdPartyPaymentID(ctx context.Context, thirdPartyPaymentID string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.Where("third_party_payment_id = ?", thirdPartyPaymentID).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) UpdateStatusByThirdPartyPaymentID(ctx context.Context, thirdPartyPaymentID string, status types.PaymentStatus) error {
	return r.db.Model(&domain.Transaction{}).Where("third_party_payment_id = ?", thirdPartyPaymentID).Update("status", status).Error
}

func (r *TransactionRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	if err := r.db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
