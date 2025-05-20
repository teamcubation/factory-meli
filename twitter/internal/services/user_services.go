package services

import (
	"context"

	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

type UserServiceImpl struct {
	repository postgres.UserRepository
}

func NewUserService(repo postgres.UserRepository) UserServiceImpl {
	return UserServiceImpl{
		repository: repo,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repository.Save(ctx, postgres.UserSaveParams{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) ListUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repository.FindByEmail(ctx, postgres.UserFindByEmailParams{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
