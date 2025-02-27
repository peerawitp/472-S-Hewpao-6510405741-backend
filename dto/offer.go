package dto

import "time"

type CreateOfferDTO struct {
	ProductRequestID uint      `json:"product_request_id"`
	OfferDate        time.Time `json:"offer_date"`
}
