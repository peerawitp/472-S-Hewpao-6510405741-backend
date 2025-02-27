package dto


type CreateTransactionRequestDTO struct {
	UserID   string `json:"user_id"`
	Amount   float64    `json:"amount"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
}