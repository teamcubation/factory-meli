package mappers

import (
	"Yuri/twitter/core/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userMapper struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	FollowUsers []string           `bson:"follow_users"`
}

func UserToMongo(user *models.User) (*userMapper, error) {
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

	return &userMapper{
		ID:          oid,
		Name:        user.Name,
		FollowUsers: user.FollowUsers,
	}, nil
}

func IdToObjectId(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return oid, nil
}

func MongoToUser(user *userMapper) *models.User {
	return &models.User{
		ID:          user.ID.Hex(),
		Name:        user.Name,
		FollowUsers: user.FollowUsers,
	}
}
