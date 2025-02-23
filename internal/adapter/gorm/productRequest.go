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

func (pr *ProductRequestGormRepo) FindByID(id int) (*domain.ProductRequest, error) {
	var productRequest domain.ProductRequest
	result := pr.db.Preload("User").Preload("Offers").First(&productRequest, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &productRequest, nil
}

func (pr *ProductRequestGormRepo) FindByUserID(id string) ([]domain.ProductRequest, error) {
	var productRequests []domain.ProductRequest
	result := pr.db.Where("user_id = ?", id).Find(&productRequests)
	if result.Error != nil {
		return nil, result.Error
	}
	return productRequests, nil
}

func (pr *ProductRequestGormRepo) FindPaginatedProductRequests(page, limit int) ([]domain.ProductRequest, int64, error) {
	var productRequests []domain.ProductRequest
	var total int64
	result := pr.db.Offset((page - 1) * limit).Limit(limit).Find(&productRequests)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	pr.db.Model(&domain.ProductRequest{}).Count(&total)
	return productRequests, total, nil
}
