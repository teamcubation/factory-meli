package services

import (
	"Yuri/twitter/core/models"
)

type IUserService interface {
	CreateUser(user *models.User) error
	FollowUser(idUser, id string) error
	UnfollowUser(idUser, id string) error
}
