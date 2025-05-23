package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	MAX_TWEET_LEN = 280
)

type User struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Tweets    []Tweet    `json:"tweets"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

type Follow struct {
	FollowerID uuid.UUID `gorm:"type:uuid;primaryKey" json:"follower_id"`
	FollowedID uuid.UUID `gorm:"type:uuid;primaryKey" json:"followed_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Followed User `gorm:"foreignKey:FollowedID;references:ID"`
	Follower User `gorm:"foreignKey:FollowerID;references:ID"`
}

type Tweet struct {
	ID        uuid.UUID  `json:"id"`
	Content   string     `json:"content"`
	UserID    uuid.UUID  `json:"user_id"`
	User      User       `json:"user"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
