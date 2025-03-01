package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hewpao/hewpao-backend/types"
	"gorm.io/gorm"
)

type Transaction struct {
	ID                  string              `gorm:"type:uuid;primaryKey"`
	UserID              string              `gorm:"index;not null"`
	Amount              float64             `gorm:"not null"`
	Currency            string              `gorm:"size:3;not null;check:currency IN ('THB')"`
	Type                string              `gorm:"size:10;not null"`
	ThirdPartyGateway   string              `gorm:"size:50;not null"`
	ThirdPartyPaymentID *string             `gorm:"size:255"`
	ProductRequestID    *uint               `gorm:"index"`
	Status              types.PaymentStatus `gorm:"type:varchar(20);default:'PENDING'"`
	CreatedAt           time.Time           `gorm:"autoCreateTime"`
	UpdatedAt           time.Time           `gorm:"autoUpdateTime"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
