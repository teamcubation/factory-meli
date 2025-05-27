package models

import "time"

type Tweet struct {
	Id        string
	UserId    string
	Content   string
	CreatedAt time.Time
}

type Timeline struct {
	TweetId   string
	UserId    string
	Name      string
	Content   string
	CreatedAt time.Time
}
