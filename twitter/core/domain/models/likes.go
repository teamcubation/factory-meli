package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	UserID    uuid.UUID
	TweetID   uuid.UUID
}
