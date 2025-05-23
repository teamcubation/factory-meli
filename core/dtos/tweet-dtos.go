package dtos

import "github.com/google/uuid"

type CreateTweetRequestDTO struct {
	UserID  uuid.UUID `json:"user_id"`
	Content string    `json:"content"`
}
