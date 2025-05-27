package models

import "time"

type Follow struct {
	FollowedId        string
	FollowerId      string
	CreatedAt time.Time
}
