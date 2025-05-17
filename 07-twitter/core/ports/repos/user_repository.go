package repos

import (
	tweetModel "07-twitter/core/models"
	// "github.com/google/uuid"
)

// UserRepository define as operações de persistência e relacionamento de usuários.
type UserRepository interface {
	// Save() armazena um usuário no repositório.
	Save(user *tweetModel.User) error

	// // Update() atualiza um usuário por seu id e retorna-o atualizado
	// Update(id uuid.UUID, updatedUser *tweetModel.User) (*tweetModel.User, error)

	// // FindByID() retorna um usuário pelo seu UUID.
	// FindByID(id uuid.UUID) (*tweetModel.User, error)

	// // Follow() faz com que o usuário com id siga o usuário followID.
	// Follow(id uuid.UUID, followID uuid.UUID) error

	// // Unfollow() faz com que o usuário com id pare de seguir o usuário followID.
	// Unfollow(id uuid.UUID, followID uuid.UUID) error

	// // ListFollowing() retorna todos os usuários que o usuário com id está seguindo.
	// ListFollowing(id uuid.UUID) ([]*tweetModel.User, error)
}
