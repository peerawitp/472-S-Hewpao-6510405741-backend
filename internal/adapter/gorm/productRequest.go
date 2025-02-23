package gorm

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
	"gorm.io/gorm"
)

type ProductRequestGormRepo struct {
	db *gorm.DB
}

func NewProductRequestGormRepo(db *gorm.DB) repository.ProductRequestRepository {
	return &ProductRequestGormRepo{db: db}
}

func (pr *ProductRequestGormRepo) Create(productRequest *domain.ProductRequest) error {
	result := pr.db.Create(&productRequest)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pr *ProductRequestGormRepo) GetDetailByID(id int) (*domain.ProductRequest, error) {
	var productRequest domain.ProductRequest
	result := pr.db.Preload("User").Preload("Offers").First(&productRequest, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &productRequest, nil
}
