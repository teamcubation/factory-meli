package models

type Tweet struct {
	ID          int
	Title       string
	Description string
	Author      User
}
