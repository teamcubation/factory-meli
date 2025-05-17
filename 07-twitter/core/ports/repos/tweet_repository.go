package repos

import (
	"07-twitter/core/models"
	"github.com/google/uuid"
)

// TweetRepository define as operações de persistência para tweets.
type TweetRepository interface {
	// Save() armazena um novo tweet a um id de usuário.
	Save(tweet *models.Tweet, userID uuid.UUID) error

	// Update() atualiza um tweet já existente.
	Update(id uuid.UUID, updatedTweet *models.Tweet) error

	// Delete() remove um tweet pelo seu ID.
	Delete(id uuid.UUID) error

	// FindAllByUserID() retorna todos os tweets de um usuário específico.
	FindAllByUserID(userID uuid.UUID) ([]*models.Tweet, error)

	// FindTimelineByUserIDs() retorna os tweets dos usuários seguidos, ordenados do mais recente para o mais antigo, com paginação.
	FindTimelineByUserIDs(userIDs []uuid.UUID, limit, offset int) ([]*models.Tweet, error)
}
