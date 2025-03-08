package domain

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserID    string
	ChatID    uint 
	Content   string

}