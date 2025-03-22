package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
)

var allowedCurrencies = map[string]bool{
	"THB": true,
}

type TransactionUseCase interface {
	CreateTransaction(userID string, amount float64, currency string) (*domain.Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (*domain.Transaction, error)
	GetTransactionByThirdPartyPaymentID(ctx context.Context, thirdPartyPaymentID string) (*domain.Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error)
}

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(tr repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: tr}
}

func (u *TransactionService) CreateTransaction(userID string, amount float64, currency string) (*domain.Transaction, error) {
	if _, exists := allowedCurrencies[currency]; !exists {
		return nil, errors.New("unsupported currency")
	}

	transaction := &domain.Transaction{
		UserID:    userID,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: time.Now(),
	}

	err := u.repo.Store(context.Background(), transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (u *TransactionService) GetTransactionByID(ctx context.Context, id string) (*domain.Transaction, error) {
	transaction, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (u *TransactionService) GetTransactionByThirdPartyPaymentID(ctx context.Context, thirdPartyPaymentID string) (*domain.Transaction, error) {
	transaction, err := u.repo.FindByThirdPartyPaymentID(ctx, thirdPartyPaymentID)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (u *TransactionService) GetTransactionsByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error) {
	transactions, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
