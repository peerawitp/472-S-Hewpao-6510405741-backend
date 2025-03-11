package usecase

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
)

type TravelerPayoutAccountUsecase interface {
	CreateTravelerPayoutAccount(userID string, req *dto.CreateTravelerPayoutAccountRequestDTO) error
	GetAccountsByUserID(userID string) (*dto.GetTravelerPayoutAccountResponseDTO, error)
	GetAllAvailableBank() (*dto.GetAllAvailableBankResponseDTO, error)
}

type travelerPayoutAccountService struct {
	ctx      context.Context
	repo     repository.TravelerPayoutAccountRepository
	bankRepo repository.BankRepository
}

func NewTravelerPayoutAccountService(ctx context.Context, repo repository.TravelerPayoutAccountRepository, bankRepo repository.BankRepository) TravelerPayoutAccountUsecase {
	return &travelerPayoutAccountService{
		ctx:      ctx,
		repo:     repo,
		bankRepo: bankRepo,
	}
}

func (t *travelerPayoutAccountService) CreateTravelerPayoutAccount(userID string, req *dto.CreateTravelerPayoutAccountRequestDTO) error {
	bank, err := t.bankRepo.GetBySwift(req.BankSwift)
	if err != nil {
		return err
	}

	account := &domain.TravelerPayoutAccount{
		UserID:        userID,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		BankSwift:     bank.SwiftCode,
	}

	storeErr := t.repo.Store(t.ctx, account)
	if storeErr != nil {
		return storeErr
	}

	return nil
}

func (t *travelerPayoutAccountService) GetAccountsByUserID(userID string) (*dto.GetTravelerPayoutAccountResponseDTO, error) {
	accounts, err := t.repo.FindByUserID(t.ctx, userID)
	if err != nil {
		return nil, err
	}

	res := &dto.GetTravelerPayoutAccountResponseDTO{
		Accounts: accounts,
	}

	return res, nil
}

func (t *travelerPayoutAccountService) GetAllAvailableBank() (*dto.GetAllAvailableBankResponseDTO, error) {
	banks, err := t.bankRepo.GetAll()
	if err != nil {
		return nil, err
	}

	res := &dto.GetAllAvailableBankResponseDTO{
		Banks: banks,
	}

	return res, nil
}
