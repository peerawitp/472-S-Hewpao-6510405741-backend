package dto

import (
	"time"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/lib/pq"
)

type UpdateProductRequestDTO struct {
	Name     string         `json:"name"`
	Desc     string         `json:"desc"`
	Quantity uint           `json:"quantity" validate:"gt=0"`
	Category types.Category `json:"category" validate:"category"`

	SelectedOfferID uint `json:"selected_offer_id"`
}

type UpdateProductRequestStatusDTO struct {
	DeliveryStatus types.DeliveryStatus `json:"delivery_status" validate:"required,delivery-status"`
	NotifyProvider string               `json:"notify_provider" validate:"required"`
}

type CreateProductRequestRequestDTO struct {
	Name         string         `json:"name" form:"name" validate:"required"`
	Desc         string         `json:"desc" form:"desc" validate:"required"`
	Budget       float64        `json:"budget" form:"budget" validate:"required,gt=0"`
	Quantity     uint           `json:"quantity" form:"quantity" validate:"required,gt=0"`
	Category     types.Category `json:"category" form:"category" validate:"required,category"`
	From         string         `json:"from" form:"from" validate:"required"`
	To           string         `json:"to" form:"to" validate:"required"`
	CheckService bool           `json:"check_service form:"check_service" validate:"required"`
}

type CreateProductRequestResponseDTO struct {
	Name     string         `json:"name"`
	Desc     string         `json:"desc"`
	Images   pq.StringArray `json:"images"`
	Budget   float64        `json:"budget"`
	Quantity uint           `json:"quantity"`
	Category types.Category `json:"category"`

	UserID         *string              `json:"userID"`
	From           string               `json:"deliver_from"`
	To             string               `json:"deliver_to"`
	CheckService   bool                 `json:"check_service"`
	DeliveryStatus types.DeliveryStatus `json:"delivery_status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type DetailOfProductRequestResponseDTO struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	Desc     string         `json:"desc"`
	Images   pq.StringArray `json:"images"`
	Budget   float64        `json:"budget"`
	Quantity uint           `json:"quantity"`
	Category types.Category `json:"category"`

	UserID *string        `json:"userID"`
	Offers []domain.Offer `json:"offers"`

	SelectedOfferID *uint `json:"selected_offer_id"`

	Transactions   []domain.Transaction `json:"transactions"`
	DeliveryStatus types.DeliveryStatus `json:"delivery_status"`

	From         string `json:"deliver_from"`
	To           string `json:"deliver_to"`
	CheckService bool   `json:"check_service"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PaginationGetProductRequestRequestDTO struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type PaginationGetProductRequestResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}
