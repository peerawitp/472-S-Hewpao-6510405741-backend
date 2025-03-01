package dto

import (
	"github.com/hewpao/hewpao-backend/types"
)

type NotificationDTO struct {
	ToID    string `json:"to_id"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type NotificationDataDTO struct {
	RecipientName string
	CompanyName   string
	ProductID     uint
	ProductStatus types.DeliveryStatus
	SupportEmail  string
	Year          int
}
