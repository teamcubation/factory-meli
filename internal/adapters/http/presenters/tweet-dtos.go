package presenters

import (
	"time"

	"github.com/google/uuid"
)

type TweetPresenter struct {
	Content   string        `json:"content"`
	UserID    uuid.UUID     `json:"user_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at"`
	User      UserPresenter `json:"user"`
}
