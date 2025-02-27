package repository

import (
	"github.com/hewpao/hewpao-backend/domain"
)

type ProductRequestRepository interface {
	Create(productRequest *domain.ProductRequest) error
	FindByID(id int) (*domain.ProductRequest, error)
	FindByUserID(id string) ([]domain.ProductRequest, error)
	FindPaginatedProductRequests(page, limit int) ([]domain.ProductRequest, int64, error)
	Update(productRequest *domain.ProductRequest) error
}
