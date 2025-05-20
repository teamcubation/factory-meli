package models

type Tweet struct {
	ID          string
	Title       string
	Description string
	Author      User
}
