package service

import (
	"time"
	"twitter-api-clone/core/models"
	"twitter-api-clone/core/ports"

	"github.com/google/uuid"
)

type UserServiceImpl struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserService) ports.UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (u *UserServiceImpl) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.Id = uuid.NewString()

	if err := u.repo.Create(user); err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) Delete(userId string) error {
	if err := u.repo.Delete(userId); err != nil {
		return err
	}

	return nil
}
