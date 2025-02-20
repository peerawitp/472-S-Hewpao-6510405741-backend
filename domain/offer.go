package domain

import (
	"gorm.io/gorm"
)

type Offer struct {
	gorm.Model
	ProductRequestID *uint
	ProductRequest   *ProductRequest
	UserID           string `gorm:"not null"`
	User             *User
}
