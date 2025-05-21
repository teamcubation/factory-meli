package services

import (
	model "07-twitter/core/models"
	repo "07-twitter/core/ports/repos"
	service "07-twitter/core/ports/services"
	utils "07-twitter/internal/utils"

	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("usuário não encontrado")

// userServiceImpl é uma implementação de UserService com o UserRepository.
type userServiceImpl struct {
	repo repo.UserRepository
}

// NewUserService() cria uma instância singleton de UserService.
func NewUserService(repo repo.UserRepository) service.UserService {
	return &userServiceImpl{repo: repo}
}

var validate = validator.New()

// CreateUser() cria um novo usuário.
func (s *userServiceImpl) CreateUser(username string) (*model.User, error) {
	newUser := model.User{
		ID:        uuid.New(),
		Username:  username,
		Following: []uuid.UUID{},
	}

	if err := validate.Struct(newUser); err != nil {
		return nil, errors.New("[CreateUser() - Validate User struct] - " + err.Error())
	}

	if err := s.repo.Save(&newUser); err != nil {
		return nil, errors.New("[CreateUser() - repo.Save()] - " + err.Error())
	}

	return &newUser, nil
}

// GetUserById() retorna um usuário pelo seu ID
func (s *userServiceImpl) GetUserById(id string) (*model.User, error) {
	idUuid, err := utils.ParseStringToUuid(id)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FindById(idUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, errors.New("[s.repo.FindById()] - " + err.Error())
	}

	return user, nil
}

func (s *userServiceImpl) Follow(userId, followingId uuid.UUID) error {
	err := s.repo.Follow(userId, followingId)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) Unfollow(userId, followingId uuid.UUID) error {
	err := s.repo.Unfollow(userId, followingId)
	if err != nil {
		return err
	}

	return nil
}
