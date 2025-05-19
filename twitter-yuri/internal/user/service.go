package user

import (
	"Yuri/twitter/core/ports/repositories"
	"Yuri/twitter/core/ports/services"
	"Yuri/twitter/dtos/request"
)

type ServiceImpl struct {
	repo   repositories.IUserRepository
	nextID int
}

func NewUserService(repo repositories.IUserRepository) services.IUserService {
	return &ServiceImpl{
		repo:   repo,
		nextID: 1,
	}
}

func (s *ServiceImpl) CreateUser(*request.CreateUserRequestDTO) error {
	return nil
}

func (s *ServiceImpl) FollowUser(idUser int, id int) error {
	panic("unimplemented")
}

func (s *ServiceImpl) UnfollowUser(idUser int, id int) error {
	panic("unimplemented")
}
