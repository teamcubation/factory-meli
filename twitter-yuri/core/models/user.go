package models

type User struct {
	ID          int
	Name        string
	FollowUsers []int
}
