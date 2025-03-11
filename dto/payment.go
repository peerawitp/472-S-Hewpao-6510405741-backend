package dto

type CreatePaymentResponseDTO struct {
	PaymentID  string `json:"payment_id"`
	PaymentURL string `json:"payment_url"`
	CreatedAt  int64  `json:"created_at"`
	ExpiredAt  int64  `json:"expired_at"`
}
