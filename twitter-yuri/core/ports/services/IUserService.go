package services

import (
	"Yuri/twitter/core/models"
)

type IUserService interface {
	CreateUser(user *models.User) error
	FollowUser(idUser, id int) error
	UnfollowUser(idUser, id int) error
}
