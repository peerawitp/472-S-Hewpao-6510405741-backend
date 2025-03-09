package dto

import "github.com/hewpao/hewpao-backend/domain"

type CreateTravelerPayoutAccountRequestDTO struct {
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	BankSwift     string `json:"bank_swift" validate:"required"`
}

type GetTravelerPayoutAccountResponseDTO struct {
	Accounts []domain.TravelerPayoutAccount `json:"accounts"`
}

type GetAllAvailableBankResponseDTO struct {
	Banks []domain.Bank `json:"banks"`
}
