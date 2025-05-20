package models

type User struct {
	ID          string `bson:"-" json:"id"`
	Name        string
	FollowUsers []string
}
