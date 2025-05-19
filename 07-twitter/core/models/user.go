package models

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID   `json:"id" validate:"required,gt=0"`
	Username  string      `json:"username" validate:"required,min=3,max=50"`
	Following []uuid.UUID `json:"following"`
}
