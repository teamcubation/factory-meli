package models

import (
	"github.com/google/uuid"
	"time"
)

type Tweet struct {
	ID        uuid.UUID `json:"id" validate:"required,gt=0"`
	UserID    uuid.UUID `json:"user_id" validate:"required,gt=0"`
	Content   string    `json:"content" validate:"required,max=280"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}
