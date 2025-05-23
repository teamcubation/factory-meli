package services

import (
	"07-twitter/core/models"
	repo "07-twitter/core/ports/repos"
	"07-twitter/core/ports/services"
	"errors"

	"github.com/google/uuid"
)

// tweetServiceImpl implementa as operações de serviço para tweets.
type tweetServiceImpl struct {
	tweetRepo repo.TweetRepository
	userRepo  repo.UserRepository
}

// NewTweetService cria uma nova instância de TweetService.
func NewTweetService(tweetRepo repo.TweetRepository, userRepo repo.UserRepository) services.TweetService {
	return &tweetServiceImpl{
		tweetRepo: tweetRepo,
		userRepo:  userRepo,
	}
}

// CreateTweet() cria um novo tweet para o usuário informado.
func (s *tweetServiceImpl) CreateTweet(content string, userID uuid.UUID) (*models.Tweet, error) {
	if content == "" {
		return nil, errors.New("conteúdo do tweet não pode ser vazio")
	}

	if userID == uuid.Nil {
		return nil, errors.New("userID é obrigatório")
	}

	tweet := &models.Tweet{
		ID:      uuid.New(),
		UserID:  userID,
		Content: content,
	}

	return s.tweetRepo.Save(tweet)
}

// GetTweetByUserId() retorna todos os tweets de um usuário específico, ordenados do mais recente para o mais antigo.
func (s *tweetServiceImpl) GetTweetByUserId(userID uuid.UUID) ([]*models.Tweet, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID é obrigatório")
	}

	return s.tweetRepo.FindAllByUserID(userID)
}

const LIMIT_TIMELINE_TWEET int = 10

// GetTimelineByUserId() retorna os tweets dos usuários seguidos pelo usuário informado, ordenados do mais recente para o mais antigo, com paginação.
func (s *tweetServiceImpl) GetTimelineByUserId(userID uuid.UUID, limit, offset int) ([]*models.Tweet, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID é obrigatório")
	}

	if limit <= 0 {
		limit = LIMIT_TIMELINE_TWEET
	}
	if offset < 0 {
		offset = 0
	}

	followingIDs, err := s.userRepo.FindFollowingIDsByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(followingIDs) == 0 {
		return []*models.Tweet{}, nil
	}

	tweets, err := s.tweetRepo.FindTimelineByUsersIds(followingIDs, limit, offset)
	if err != nil {
		return nil, err
	}
	return tweets, nil
}
