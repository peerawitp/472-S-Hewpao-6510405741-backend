package dto

type CreateMessageRequestDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ChatID    uint `json:"chat_id"`
	Content   string `json:"content"`
}