package repository

import (
	"github.com/hewpao/hewpao-backend/domain"
)

type ProductRequestRepository interface {
	Create(productRequest *domain.ProductRequest) error
	GetDetailByID(id int) (*domain.ProductRequest, error)
}
