package gorm

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/repository"
	"gorm.io/gorm"
)

type BankGormRepo struct {
	db *gorm.DB
}

func NewBankGormRepo(db *gorm.DB) repository.BankRepository {
	return &BankGormRepo{db: db}
}

func (b *BankGormRepo) GetBySwift(swift string) (*domain.Bank, error) {
	bank := domain.Bank{}
	err := b.db.First(&bank, "swift_code = ?", swift).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.ErrBankNotFound
		}
		return nil, err
	}
	return &bank, nil
}

func (b *BankGormRepo) GetAll() ([]domain.Bank, error) {
	banks := []domain.Bank{}
	err := b.db.Find(&banks).Error
	if err != nil {
		return nil, err
	}
	return banks, nil
}
