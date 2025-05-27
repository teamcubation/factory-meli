package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	Post      string
	CreatorID uuid.UUID
	ParentID  uuid.NullUUID
}
