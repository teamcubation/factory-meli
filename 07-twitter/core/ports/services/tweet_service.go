package services

import (
	tweetModel "07-twitter/core/models"

	"github.com/google/uuid"
)

// TweetService define as operações de persistência para tweets.
type TweetService interface {
	// CreateTweet() cria um novo tweet associando a um id de usuário.
	CreateTweet(content string, userID uuid.UUID) (*tweetModel.Tweet, error)

	// GetTweetByUserId() retorna todos os tweets de um usuário específico.
	GetTweetByUserId(userID uuid.UUID) ([]*tweetModel.Tweet, error)

	// FindTimelineByUserId() retorna os tweets dos usuários seguidos, ordenados do mais recente para o mais antigo, com paginação.
	GetTimelineByUserId(userID uuid.UUID, limit, offset int) ([]*tweetModel.Tweet, error)
}
