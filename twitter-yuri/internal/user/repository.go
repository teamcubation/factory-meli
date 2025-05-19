package user

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/db"
	"context"
)

type UserRepositoryImpl struct{}

func NewUserRepo() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (u *UserRepositoryImpl) SaveUser(user *models.User) error {
	collection, err := db.ConnectToCollection("users")
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepositoryImpl) UpdateUser(user *models.User) error {
	collection, err := db.ConnectToCollection("users")
	if err != nil {
		return err
	}
	_, err = collection.UpdateByID(context.TODO(), user.ID, user)
	if err != nil {
		return err
	}
	return nil
}
