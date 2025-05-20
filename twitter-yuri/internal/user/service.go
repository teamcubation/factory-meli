package user

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/repositories"
	"Yuri/twitter/core/ports/services"
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

func (s *ServiceImpl) CreateUser(user *models.User) error {
	err := s.repo.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) FollowUser(idUser int, id int) error {
	panic("unimplemented")
}

func (s *ServiceImpl) UnfollowUser(idUser int, id int) error {
	panic("unimplemented")
}
