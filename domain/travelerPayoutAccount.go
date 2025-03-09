package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TravelerPayoutAccount struct {
	ID            string `gorm:"primaryKey"`
	UserID        string `gorm:"not null"`
	AccountNumber string `gorm:"size:50;not null"`
	AccountName   string `gorm:"size:50;not null"`
	BankSwift     string `gorm:"size:100;not null"`
	Bank          Bank   `gorm:"foreignKey:BankSwift"`
}

func (t *TravelerPayoutAccount) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
