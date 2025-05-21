package models

import "time"

type Tweet struct {
	ID          string
	Title       string
	Description string
	AuthorId    string
	CreatedAt   time.Time
}
