package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction(t *testing.T) {
	mockTransactionRepo := new(mock_repos.MockTransactionRepository)
	transactionService := usecase.NewTransactionService(mockTransactionRepo)
	ctx := context.Background()

	// Mock transaction data
	expectedTransaction := &domain.Transaction{
		UserID:   "1",
		Amount:   123.0,
		Currency: "THB",
	}

	t.Run("Success", func(t *testing.T) {
		mockTransactionRepo.ExpectedCalls = nil

		// Mock repository response
		mockTransactionRepo.On("Store", ctx, mock.MatchedBy(func(t *domain.Transaction) bool {
			// Check that the fields match
			return t.UserID == expectedTransaction.UserID && t.Amount == expectedTransaction.Amount && t.Currency == expectedTransaction.Currency
		})).Return(nil).Once()

		// Call the function
		transaction, err := transactionService.CreateTransaction(expectedTransaction.UserID, expectedTransaction.Amount, expectedTransaction.Currency)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, expectedTransaction.UserID, transaction.UserID)
		assert.Equal(t, expectedTransaction.Amount, transaction.Amount)
		assert.Equal(t, expectedTransaction.Currency, transaction.Currency)

		// Verify expectations
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Error_UnsupportedCurrency", func(t *testing.T) {
		mockTransactionRepo.ExpectedCalls = nil

		unsupportedCurrency := "alkwjdS"

		_, err := transactionService.CreateTransaction(expectedTransaction.UserID, expectedTransaction.Amount, unsupportedCurrency)

		assert.Error(t, err)
		assert.Equal(t, errors.New("unsupported currency"), err)

		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Error_StoreTransaction", func(t *testing.T) {
		mockTransactionRepo.ExpectedCalls = nil

		// Mock repository response
		mockTransactionRepo.On("Store", ctx, mock.Anything).Return(errors.New("database error")).Once()

		// Call the function
		transaction, err := transactionService.CreateTransaction(expectedTransaction.UserID, expectedTransaction.Amount, expectedTransaction.Currency)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.EqualError(t, err, "database error")

		// Verify expectations
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestGetTransactionByID(t *testing.T) {
	mockTransactionRepo := new(mock_repos.MockTransactionRepository)
	transactionService := usecase.NewTransactionService(mockTransactionRepo)
	ctx := context.Background()

	// Mock transaction data
	expectedTransaction := &domain.Transaction{
		ID:       "txn_123",
		UserID:   "1",
		Amount:   123.0,
		Currency: "THB",
	}

	t.Run("Success", func(t *testing.T) {
		mockTransactionRepo.ExpectedCalls = nil

		// Mock repository response
		mockTransactionRepo.On("FindByID", ctx, expectedTransaction.ID).Return(expectedTransaction, nil).Once()

		// Call the function
		transaction, err := transactionService.GetTransactionByID(ctx, expectedTransaction.ID)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, expectedTransaction.ID, transaction.ID)
		assert.Equal(t, expectedTransaction.UserID, transaction.UserID)
		assert.Equal(t, expectedTransaction.Amount, transaction.Amount)
		assert.Equal(t, expectedTransaction.Currency, transaction.Currency)

		// Verify expectations
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Error_FindByID", func(t *testing.T) {
		mockTransactionRepo.ExpectedCalls = nil

		// Mock repository response for error case
		mockTransactionRepo.On("FindByID", ctx, expectedTransaction.ID).Return((*domain.Transaction)(nil), errors.New("not found")).Once()

		// Call the function
		transaction, err := transactionService.GetTransactionByID(ctx, expectedTransaction.ID)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.EqualError(t, err, "not found")

		// Verify expectations
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestGetTransactionByThirdPartyPaymentID(t *testing.T) {
	mockTransactionRepo := new(mock_repos.MockTransactionRepository)
	transactionService := usecase.NewTransactionService(mockTransactionRepo)
	ctx := context.Background()

	// Mock transaction data
	third := "thirdParty_123"
	expectedTransaction := &domain.Transaction{
		ID:                  "txn_123",
		UserID:              "1",
		Amount:              123.0,
		Currency:            "THB",
		ThirdPartyPaymentID: &third,
	}

	t.Run("Success", func(t *testing.T) {
		// Mock repository response
		mockTransactionRepo.On("FindByThirdPartyPaymentID", ctx, *expectedTransaction.ThirdPartyPaymentID).
			Return(expectedTransaction, nil).Once()

		// Call the function
		transaction, err := transactionService.GetTransactionByThirdPartyPaymentID(ctx, *expectedTransaction.ThirdPartyPaymentID)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, expectedTransaction.ID, transaction.ID)
		assert.Equal(t, *expectedTransaction.ThirdPartyPaymentID, *transaction.ThirdPartyPaymentID)

		// Verify expectations
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Error_TransactionNotFound", func(t *testing.T) {
		// Mock repository response for error case
		mockTransactionRepo.On("FindByThirdPartyPaymentID", ctx, "nonexistent_id").
			Return(&domain.Transaction{}, errors.New("transaction not found")).Once()

		// Call the function
		transaction, err := transactionService.GetTransactionByThirdPartyPaymentID(ctx, "nonexistent_id")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.EqualError(t, err, "transaction not found")

		// Verify expectations
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestGetTransactionsByUserID(t *testing.T) {
	mockTransactionRepo := new(mock_repos.MockTransactionRepository)
	transactionService := usecase.NewTransactionService(mockTransactionRepo)
	ctx := context.Background()

	// Mock transaction data
	expectedTransactions := []*domain.Transaction{
		{
			ID:       "txn_123",
			UserID:   "1",
			Amount:   100.0,
			Currency: "THB",
		},
		{
			ID:       "txn_124",
			UserID:   "1",
			Amount:   200.0,
			Currency: "USD",
		},
	}

	t.Run("Success", func(t *testing.T) {
		// Mock repository response
		mockTransactionRepo.On("FindByUserID", ctx, "1").
			Return(expectedTransactions, nil).Once()

		// Call the function
		transactions, err := transactionService.GetTransactionsByUserID(ctx, "1")

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, transactions)
		assert.Len(t, transactions, 2)
		assert.Equal(t, expectedTransactions[0].ID, transactions[0].ID)
		assert.Equal(t, expectedTransactions[1].ID, transactions[1].ID)

		// Verify expectations
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Error_NoTransactionsFound", func(t *testing.T) {
		// Mock repository response for error case
		mockTransactionRepo.On("FindByUserID", ctx, "nonexistent_user").
			Return([]*domain.Transaction{}, errors.New("no transactions found")).Once()

		// Call the function
		transactions, err := transactionService.GetTransactionsByUserID(ctx, "nonexistent_user")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, transactions)
		assert.EqualError(t, err, "no transactions found")

		// Verify expectations
		mockTransactionRepo.AssertExpectations(t)
	})
}
