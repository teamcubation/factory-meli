package repos

import (
	model "07-twitter/core/models"

	"github.com/google/uuid"
)

// UserRepository define as operações de persistência e relacionamento de usuários.
type UserRepository interface {
	// Save() armazena um usuário no repositório.
	Save(user *model.User) error

	// FindById() buscar um usuário pelo seu ID.
	FindById(id uuid.UUID) (*model.User, error)

	// Follow() faz com que o usuário com userId siga o usuário followId.
	Follow(userId, followingId uuid.UUID) error

	// Unfollow() faz com que o usuário com userId deixe de seguir o usuário followId.
	Unfollow(userId, followingId uuid.UUID) error
}
