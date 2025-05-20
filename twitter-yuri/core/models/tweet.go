package models

type Tweet struct {
	ID          string `bson:"-" json:"id"`
	Title       string
	Description string
	Author      User
}
