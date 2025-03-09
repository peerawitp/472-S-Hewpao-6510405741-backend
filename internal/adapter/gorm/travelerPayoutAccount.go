package gorm

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/repository"
	"gorm.io/gorm"
)

type TravelerPayoutAccountGormRepository struct {
	db *gorm.DB
}

func NewTravelerPayoutAccountGormRepository(db *gorm.DB) repository.TravelerPayoutAccountRepository {
	return &TravelerPayoutAccountGormRepository{db: db}
}

func (t *TravelerPayoutAccountGormRepository) Store(ctx context.Context, travelerPayoutAccount *domain.TravelerPayoutAccount) error {
	exist := t.db.Where("account_number = ? AND bank_swift = ?", travelerPayoutAccount.AccountNumber, travelerPayoutAccount.BankSwift).First(&domain.TravelerPayoutAccount{})
	if exist.RowsAffected == 1 {
		return exception.ErrDuplicateTravelerPayoutAccount
	}
	result := t.db.Create(travelerPayoutAccount)
	return result.Error
}

func (t *TravelerPayoutAccountGormRepository) FindByUserID(ctx context.Context, userID string) ([]domain.TravelerPayoutAccount, error) {
	var travelerPayoutAccounts []domain.TravelerPayoutAccount
	result := t.db.Where("user_id = ?", userID).Preload("Bank").Find(&travelerPayoutAccounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return travelerPayoutAccounts, nil
}
