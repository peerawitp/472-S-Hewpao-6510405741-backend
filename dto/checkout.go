package dto

type CheckoutRequestDTO struct {
	ProductRequestID uint   `json:"product_request_id" validate:"required"`
	PaymentGateway   string `json:"payment_gateway" validate:"required"`
}

type CheckoutResponseDTO struct {
	Payment *CreatePaymentResponseDTO `json:"payment"`
}
