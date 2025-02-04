package domain

import (
	"time"
)

type User struct {
	ID          string  `gorm:"primaryKey"`
	Name        string  `gorm:"size:50;not null"`
	MiddleName  *string `gorm:"size:50"`
	Surname     string  `gorm:"size:50;not null"`
	Email       string  `gorm:"unique;not null"`
	Password    *string `gorm:"size:255" json:"-"`
	PhoneNumber *string `gorm:"unique;size:15"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	Accounts    []Account `gorm:"foreignKey:UserID"`
}
