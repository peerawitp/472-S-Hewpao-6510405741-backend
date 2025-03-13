package dto

import (
	"time"

	"github.com/hewpao/hewpao-backend/domain"
)

type CreateOfferDTO struct {
	ProductRequestID uint      `json:"product_request_id"`
	OfferDate        time.Time `json:"offer_date"`
}

type GetOfferDetailDTO struct {
	ID               uint                   `json:"id"`
	ProductRequestID *uint                  `json:"product_request_id"`
	ProductRequest   *domain.ProductRequest `json:"product_request"`
	UserID           string                 `json:"user_id"`
	User             *domain.User           `json:"user"`
	OfferDate        time.Time              `json:"offer_date"`
}
