package services

import (
	userModel "07-twitter/core/models"

	"github.com/google/uuid"
)

// UserService define os casos de uso relacionados a usuários.
type UserService interface {
	// CreateUser() cria um novo usuário.
	CreateUser(username string) (*userModel.User, error)

	// GetUserById() retorna um usuário pelo seu ID.
	GetUserById(id string) (*userModel.User, error)

	// Follow() faz um usuário com userId seguir outro usuário com followingId
	Follow(userId, followingId uuid.UUID) error

	// Unfollow() faz um usuário com userId deixar de seguir outro com followingId
	Unfollow(userId, followingId uuid.UUID) error
}
