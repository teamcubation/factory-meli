package models

type User struct {
	ID          string
	Name        string
	FollowUsers []string
}
