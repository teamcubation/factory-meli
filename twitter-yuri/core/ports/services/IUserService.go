package services

import "Yuri/twitter/dtos/request"

type IUserService interface {
	CreateUser(user *request.CreateUserRequestDTO) error
	FollowUser(idUser, id int) error
	UnfollowUser(idUser, id int) error
}
