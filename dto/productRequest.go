package dto

import (
	"time"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/lib/pq"
)

type CreateProductRequestRequestDTO struct {
	Name     string         `json:"name" validate:"required"`
	Desc     string         `json:"desc" validate:"required"`
	Budget   float64        `json:"budget" validate:"required,gt=0"`
	Quantity uint           `json:"quantity" validate:"required,gt=0"`
	Category types.Category `json:"category" validate:"required,category"`
}

type CreateProductRequestResponseDTO struct {
	Name     string         `json:"name"`
	Desc     string         `json:"desc"`
	Images   pq.StringArray `json:"images"`
	Budget   float64        `json:"budget"`
	Quantity uint           `json:"quantity"`
	Category types.Category `json:"category"`

	UserID *string        `json:"userID"`
	User   *domain.User   `json:"user"`
	Offers []domain.Offer `json:"offers"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type DetailOfProductRequestResponseDTO struct {
	ID       uint           `json:"id"`
	Desc     string         `json:"desc"`
	Images   pq.StringArray `json:"images"`
	Budget   float64        `json:"budget"`
	Quantity uint           `json:"quantity"`
	Category types.Category `json:"category"`

	UserID *string        `json:"userID"`
	User   *domain.User   `json:"user"`
	Offers []domain.Offer `json:"offers"`

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
