package tweet

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/out"
	"Yuri/twitter/core/ports/services"
)

type TweetServiceImpl struct {
	repo   out.ITweetRepository
	nextID int
}

func NewTweetService(repo out.ITweetRepository) services.ITweetService {
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
