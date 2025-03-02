package repository

import "github.com/hewpao/hewpao-backend/domain"

type VerificationRepository interface {
	Create(verification *domain.Verification) error
	FindByID(id uint) (*domain.Verification, error)
}
