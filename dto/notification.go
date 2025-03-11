package dto

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/types"
)

type SendNotificationDTO struct {
	ToUser  domain.User `json:"to_user"`
	Subject string      `json:"subject"`
	Content string      `json:"content"`
}

type NotificationDataDTO struct {
	RecipientName string
	CompanyName   string
	ProductID     uint
	ProductStatus types.DeliveryStatus
	SupportEmail  string
	Year          int
}
