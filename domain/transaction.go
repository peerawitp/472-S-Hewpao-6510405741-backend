package domain

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Transaction struct {
    ID        string    `gorm:"type:uuid;primaryKey"`
    UserID    string    `gorm:"index;not null"`
    Amount    float64   `gorm:"not null"`
    Currency  string    `gorm:"size:3;not null;check:currency IN ('THB')"`
    Type      string    `gorm:"size:10;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
    t.ID = uuid.New().String()
    return
}
