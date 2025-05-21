package repositories

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/infrastructure/db/mongo"
	"Yuri/twitter/infrastructure/db/mongo/mappers"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepositoryImpl struct{}

func NewUserRepo() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (u *UserRepositoryImpl) SaveUser(user *models.User) error {
	collection, err := mongo.ConnectToCollection("users")
	if err != nil {
		mongo.Disconnect(collection.Database().Client())
		return err
	}
	mongoUser, err := mappers.UserToMongo(user)
	if err != nil {
		mongo.Disconnect(collection.Database().Client())
		return err
	}

	if _, errInsert := collection.InsertOne(context.TODO(), mongoUser); errInsert != nil {
		mongo.Disconnect(collection.Database().Client())
		return errInsert
	}

	if errDisconnect := mongo.Disconnect(collection.Database().Client()); errDisconnect != nil {
		return errDisconnect
	}
	return nil
}

func (u *UserRepositoryImpl) UpdateUser(user *models.User) error {
	collection, err := mongo.ConnectToCollection("users")
	if err != nil {
		mongo.Disconnect(collection.Database().Client())
		return err
	}
	mongoUser, errMapper := mappers.UserToMongo(user)
	if errMapper != nil {
		mongo.Disconnect(collection.Database().Client())
		return errMapper
	}

	if _, errUpdate := collection.UpdateByID(context.TODO(), mongoUser.ID, bson.M{"$set": mongoUser}); errUpdate != nil {
		mongo.Disconnect(collection.Database().Client())
		return errUpdate
	}

	if errDisconnect := mongo.Disconnect(collection.Database().Client()); errDisconnect != nil {
		return errDisconnect
	}
	return nil
}

func (u *UserRepositoryImpl) GetById(id string) (*models.User, error) {
	var user models.User
	collection, err := mongo.ConnectToCollection("users")
	if err != nil {
		mongo.Disconnect(collection.Database().Client())
		return nil, err
	}

	objectID, err := mappers.IdToObjectId(id)
	if err != nil {
		mongo.Disconnect(collection.Database().Client())
		return nil, err
	}

	result := collection.FindOne(context.TODO(), bson.M{"_id": objectID})
	if result.Err() != nil {
		mongo.Disconnect(collection.Database().Client())
		return nil, result.Err()
	}

	userMapper, errMap := mappers.UserToMongo(&user)
	if errMap != nil {
		mongo.Disconnect(collection.Database().Client())
		return nil, errMap
	}

	if errDecode := result.Decode(userMapper); errDecode != nil {
		mongo.Disconnect(collection.Database().Client())
		return nil, errDecode
	}
	user = *mappers.MongoToUser(userMapper)
	if errDisconnect := mongo.Disconnect(collection.Database().Client()); errDisconnect != nil {
		return nil, errDisconnect
	}
	return &user, nil
}
