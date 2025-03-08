package domain

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Name             string
	Message string `gorm:"foreignKey:ChatID"`
}
