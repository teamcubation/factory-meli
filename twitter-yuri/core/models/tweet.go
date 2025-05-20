package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tweet struct {
	ID          primitive.ObjectID
	Title       string
	Description string
	Author      User
}
