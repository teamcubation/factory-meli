package services

import (
	userModel "07-twitter/core/models"

	// "github.com/google/uuid"
)

// UserService define os casos de uso relacionados a usuários.
type UserService interface {
	// CreateUser() cria um novo usuário.
	CreateUser(username string) (*userModel.User, error)

	// UpdateUser() atualiza um usuário por seu id e retorna-o atualizado
	// UpdateUser(id uuid.UUID, updatedUser *userModel.User) (*userModel.User, error)

	// // GetUserByID() busca e retorna um usuário pelo seu id
	// GetUserByID(id uuid.UUID) (*userModel.User, error)

	// // Follow() faz com que o usuário com id siga o usuário followID.
	// Follow(id uuid.UUID, followID uuid.UUID) error

	// // Unfollow() faz com que o usuário com id pare de seguir o usuário followID.
	// Unfollow(id uuid.UUID, followID uuid.UUID) error

	// // ListFollowers() retorna todos os seguidores do usuário com id.
	// ListFollowers(id uuid.UUID) ([]*userModel.User, error)
}
