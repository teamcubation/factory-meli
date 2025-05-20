package tweet

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/repositories"
	"Yuri/twitter/core/ports/services"
)

type TweetServiceImpl struct {
	repo   repositories.ITweetRepository
	nextID int
}

func NewTweetService(repo repositories.ITweetRepository) services.ITweetService {
	return &TweetServiceImpl{
		repo:   repo,
		nextID: 1,
	}
}

func (s *TweetServiceImpl) CreateTweet(tweet *models.Tweet) error {
	panic("unimplemented")
}

func (s *TweetServiceImpl) ListTweetsByUserID(id int) ([]*models.Tweet, error) {
	panic("unimplemented")
}
