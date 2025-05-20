package utils

import (
	"Yuri/twitter/core/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoUser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	FollowUsers []string           `bson:"followusers"`
}

func ToMongo(user *models.User) (*mongoUser, error) {
	var oid primitive.ObjectID
	var err error
	if user.ID == "" {
		oid = primitive.NewObjectID()
	} else {
		oid, err = primitive.ObjectIDFromHex(user.ID)
		if err != nil {
			return nil, err
		}
	}

	return &mongoUser{
		ID:          oid,
		Name:        user.Name,
		FollowUsers: user.FollowUsers,
	}, nil
}
