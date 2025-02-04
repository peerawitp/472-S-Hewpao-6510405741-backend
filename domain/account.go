package domain

import "time"

type Account struct {
	ID                string  `gorm:"primaryKey"`
	UserID            string  `gorm:"not null"`
	Provider          string  `gorm:"size:50;not null"`
	ProviderAccountID string  `gorm:"size:255;not null"`
	Type              string  `gorm:"size:50;not null"`
	AccessToken       *string `gorm:"size:512"`
	RefreshToken      *string `gorm:"size:512"`
	TokenType         *string `gorm:"size:50"`
	Scope             *string `gorm:"size:512"`
	ExpiredAt         *time.Time
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}
