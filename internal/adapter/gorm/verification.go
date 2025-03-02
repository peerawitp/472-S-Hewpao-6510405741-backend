package gorm

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
	"gorm.io/gorm"
)

type verificationGormRepo struct {
	db *gorm.DB
}

func NewVerificationGormRepo(db *gorm.DB) repository.VerificationRepository {
	return &verificationGormRepo{
		db: db,
	}
}

func (v *verificationGormRepo) Create(verification *domain.Verification) error {
	result := v.db.Create(&verification)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (v *verificationGormRepo) FindByID(id uint) (*domain.Verification, error) {
	var verification domain.Verification
	result := v.db.First(&verification, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &verification, nil
}
