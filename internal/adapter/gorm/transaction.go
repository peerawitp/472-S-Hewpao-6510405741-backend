package gorm

import (
	"context"
	"github.com/hewpao/hewpao-backend/domain"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Store(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

