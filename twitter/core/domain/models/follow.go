package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	FollowerID uuid.UUID
	FollowedID uuid.UUID
}
