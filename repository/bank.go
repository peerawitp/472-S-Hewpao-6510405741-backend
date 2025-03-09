package repository

import "github.com/hewpao/hewpao-backend/domain"

type BankRepository interface {
	GetBySwift(swift string) (*domain.Bank, error)
	GetAll() ([]domain.Bank, error)
}
