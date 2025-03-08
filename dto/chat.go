package dto

import "time"

type CreateChatRequestDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	OrderID uint `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
}
