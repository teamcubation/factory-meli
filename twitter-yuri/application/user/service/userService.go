package user

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/out"
	"Yuri/twitter/core/ports/services"
	"fmt"
)

type UserServiceImpl struct {
	repo out.IUserRepository
}

func NewUserService(repo out.IUserRepository) services.IUserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) CreateUser(user *models.User) error {
	err := s.repo.SaveUser(user)
	if err != nil {
		fmt.Println("Error saving user:", err)
		return err
	}
	return nil
}

func (s *UserServiceImpl) FollowUser(idUser int, id int) error {
	panic("unimplemented")
}

func (s *UserServiceImpl) UnfollowUser(idUser int, id int) error {
	panic("unimplemented")
}
