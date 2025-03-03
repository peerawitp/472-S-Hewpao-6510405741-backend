package dto

import (
	"github.com/hewpao/hewpao-backend/types"
)

type NotificationDTO struct {
	ToEmail string `json:"to_email"`
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
