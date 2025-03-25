package usecase_test

import (
	"context"
	"testing"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	"github.com/stretchr/testify/assert"
)

func TestCreateTravelerPayoutAccount(t *testing.T) {
	mockTpaRepo := new(mock_repos.MockTravelerPayoutAccountRepository)
	mockBankRepo := new(mock_repos.MockBankRepository)
	ctx := context.Background()

	travelerPayoutAccountService := usecase.NewTravelerPayoutAccountService(ctx, mockTpaRepo, mockBankRepo)

	t.Run("Success", func(t *testing.T) {
		mockTpaRepo.ExpectedCalls = nil
		mockBankRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		req := &dto.CreateTravelerPayoutAccountRequestDTO{
			AccountNumber: "1234",
			AccountName:   "my-traveler",
			BankSwift:     "bankswift1234",
		}
		account := &domain.TravelerPayoutAccount{
			UserID:        travelerID,
			AccountNumber: req.AccountNumber,
			AccountName:   req.AccountName,
			BankSwift:     "code1234",
		}

		// TODO: repo flow
		mockBankRepo.On("GetBySwift", req.BankSwift).Return(&domain.Bank{
			SwiftCode: "code1234",
		}, nil)
		mockTpaRepo.On("Store", ctx, account).Return(nil)

		err := travelerPayoutAccountService.CreateTravelerPayoutAccount(travelerID, req)

		assert.NoError(t, err)
		mockTpaRepo.AssertExpectations(t)
		mockBankRepo.AssertExpectations(t)
	})

	t.Run("Error_BankNotFound", func(t *testing.T) {
		mockTpaRepo.ExpectedCalls = nil
		mockBankRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		req := &dto.CreateTravelerPayoutAccountRequestDTO{
			AccountNumber: "1234",
			AccountName:   "my-traveler",
			BankSwift:     "bankswift1234",
		}

		// TODO: repo flow
		mockBankRepo.On("GetBySwift", req.BankSwift).Return(&domain.Bank{}, assert.AnError)

		err := travelerPayoutAccountService.CreateTravelerPayoutAccount(travelerID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockBankRepo.AssertExpectations(t)
	})
	t.Run("Error_CouldNotStore", func(t *testing.T) {
		mockTpaRepo.ExpectedCalls = nil
		mockBankRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		req := &dto.CreateTravelerPayoutAccountRequestDTO{
			AccountNumber: "1234",
			AccountName:   "my-traveler",
			BankSwift:     "bankswift1234",
		}
		account := &domain.TravelerPayoutAccount{
			UserID:        travelerID,
			AccountNumber: req.AccountNumber,
			AccountName:   req.AccountName,
			BankSwift:     "code1234",
		}

		// TODO: repo flow
		mockBankRepo.On("GetBySwift", req.BankSwift).Return(&domain.Bank{
			SwiftCode: "code1234",
		}, nil)
		mockTpaRepo.On("Store", ctx, account).Return(assert.AnError)

		err := travelerPayoutAccountService.CreateTravelerPayoutAccount(travelerID, req)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockTpaRepo.AssertExpectations(t)
		mockBankRepo.AssertExpectations(t)
	})
}

func TestGetAccountsByUserID(t *testing.T) {
	mockTpaRepo := new(mock_repos.MockTravelerPayoutAccountRepository)
	mockBankRepo := new(mock_repos.MockBankRepository)
	ctx := context.Background()

	travelerPayoutAccountService := usecase.NewTravelerPayoutAccountService(ctx, mockTpaRepo, mockBankRepo)

	t.Run("Success", func(t *testing.T) {
		mockTpaRepo.ExpectedCalls = nil
		mockBankRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		accounts := []domain.TravelerPayoutAccount{
			{},
		}

		// TODO: repo flow
		mockTpaRepo.On("FindByUserID", ctx, travelerID).Return(accounts, nil)

		_, err := travelerPayoutAccountService.GetAccountsByUserID(travelerID)

		assert.NoError(t, err)
		mockTpaRepo.AssertExpectations(t)
	})

	t.Run("Error_CouldNotGetAllAccountsForUserID", func(t *testing.T) {
		mockTpaRepo.ExpectedCalls = nil
		mockBankRepo.ExpectedCalls = nil

		travelerID := "traveler123"
		accounts := []domain.TravelerPayoutAccount{
			{},
		}

		// TODO: repo flow
		mockTpaRepo.On("FindByUserID", ctx, travelerID).Return(accounts, assert.AnError)

		_, err := travelerPayoutAccountService.GetAccountsByUserID(travelerID)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockTpaRepo.AssertExpectations(t)
	})
}

func TestGetAllAvailableBank(t *testing.T) {
	mockTpaRepo := new(mock_repos.MockTravelerPayoutAccountRepository)
	mockBankRepo := new(mock_repos.MockBankRepository)
	ctx := context.Background()

	travelerPayoutAccountService := usecase.NewTravelerPayoutAccountService(ctx, mockTpaRepo, mockBankRepo)
	t.Run("Success", func(t *testing.T) {
		mockTpaRepo.ExpectedCalls = nil
		mockBankRepo.ExpectedCalls = nil

		// TODO: repo flow
		mockBankRepo.On("GetAll").Return([]domain.Bank{}, nil)
		_, err := travelerPayoutAccountService.GetAllAvailableBank()

		assert.NoError(t, err)
		mockTpaRepo.AssertExpectations(t)
	})
	t.Run("Error_CouldNotGetAllAvaialbeBanks", func(t *testing.T) {
		mockTpaRepo.ExpectedCalls = nil
		mockBankRepo.ExpectedCalls = nil

		// TODO: repo flow
		mockBankRepo.On("GetAll").Return([]domain.Bank{}, assert.AnError)
		_, err := travelerPayoutAccountService.GetAllAvailableBank()

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockTpaRepo.AssertExpectations(t)
	})
}
