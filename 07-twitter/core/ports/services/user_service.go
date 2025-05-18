package services

import (
	userModel "07-twitter/core/models"
)

// UserService define os casos de uso relacionados a usuários.
type UserService interface {
	// CreateUser() cria um novo usuário.
	CreateUser(username string) (*userModel.User, error)

	// GetUserById() retorna um usuário pelo seu ID.
	GetUserById(id string) (*userModel.User, error)
}
