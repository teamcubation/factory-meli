package repositories

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/infrastructure/db/mongo"
	"Yuri/twitter/infrastructure/db/utils"
	"context"
	"fmt"
)

type UserRepositoryImpl struct{}

func NewUserRepo() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (u *UserRepositoryImpl) SaveUser(user *models.User) error {
	collection, err := mongo.ConnectToCollection("users")
	if err != nil {
		fmt.Println("Error to ConnectToCollection:", err)
		return err
	}
	mongoUser, err := utils.ToMongo(user)
	if err != nil {
		fmt.Println("Error to ToMongo:", err)
		return err
	}
	_, errInsert := collection.InsertOne(context.TODO(), mongoUser)
	if errInsert != nil {
		fmt.Println("Error to InsertOne:", errInsert)
		return errInsert
	}
	errDisconnect := mongo.Disconnect(collection.Database().Client())
	if errDisconnect != nil {
		fmt.Println("Error disconnecting from MongoDB:", errDisconnect)
		return errDisconnect
	}
	return nil
}

func (u *UserRepositoryImpl) UpdateUser(user *models.User) error {
	collection, err := mongo.ConnectToCollection("users")
	if err != nil {
		return err
	}
	_, err = collection.UpdateByID(context.TODO(), user.ID, user)
	if err != nil {
		return err
	}
	errDisconnect := mongo.Disconnect(collection.Database().Client())
	if errDisconnect != nil {
		fmt.Println("Error disconnecting from MongoDB:", errDisconnect)
		return errDisconnect
	}
	return nil
}
