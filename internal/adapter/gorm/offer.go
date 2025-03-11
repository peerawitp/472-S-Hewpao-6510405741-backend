package gorm

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
	"gorm.io/gorm"
)

type offerGormRepo struct {
	db *gorm.DB
}

func NewOfferGormRepo(db *gorm.DB) repository.OfferRepository {
	return &offerGormRepo{
		db: db,
	}
}

func (o *offerGormRepo) GetByID(offer *domain.Offer) error {
	result := o.db.Preload("User").Preload("ProductRequest").First(offer, offer.ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (o *offerGormRepo) Create(offer *domain.Offer) error {
	result := o.db.Create(offer)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
