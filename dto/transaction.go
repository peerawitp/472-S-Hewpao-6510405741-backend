package dto

import (
	"time"

	"github.com/hewpao/hewpao-backend/types"
)

type CreateTransactionRequestDTO struct {
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Type     string  `json:"type"`
}

type GetTransactionResponseDTO struct {
	ID                  string              `json:"id"`
	UserID              string              `json:"user_id"`
	Amount              float64             `json:"amount"`
	Currency            string              `json:"currency"`
	ThirdPartyGateway   string              `json:"third_party_gateway"`
	ThirdPartyPaymentID *string             `json:"third_party_payment_id"`
	ProductRequestID    *uint               `json:"product_request_id"`
	Status              types.PaymentStatus `json:"status"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

