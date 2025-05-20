package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

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
	// Validating email
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if email == "" {
		return nil, fmt.Errorf("email is an obligatory field")
	}

	if !emailRegex.MatchString(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Verifying if the user email already has an account associated to it
	user, err := s.repository.FindByEmail(ctx, postgres.UserFindByEmailParams{
		Email: email,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("there already exists an user with this email")
	}

	user, err = s.repository.Save(ctx, postgres.UserSaveParams{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
