package domain

import (
	"time"

	"gorm.io/gorm"
)

type Offer struct {
	gorm.Model
	ProductRequestID *uint
	ProductRequest   *ProductRequest
	UserID           string `gorm:"not null"`
	User             *User
	OfferDate        time.Time
}
