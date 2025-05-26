package tweet

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/out"
	"Yuri/twitter/core/ports/services"
	"fmt"
)

type TweetServiceImpl struct {
	repo out.ITweetRepository
}

func NewTweetService(repo out.ITweetRepository) services.ITweetService {
	return &TweetServiceImpl{
		repo: repo,
	}
}

func (s *TweetServiceImpl) CreateTweet(tweet *models.Tweet) error {
	if tweet == nil {
		return fmt.Errorf("tweet cannot be nil")
	}
	if err := s.repo.SaveTweet(tweet); err != nil {
		return err
	}
	return nil
}

func (s *TweetServiceImpl) ListTweetsByUserID(userId string, page, tweetsPerPage int) ([]*models.Tweet, error) {
	tweets, err := s.repo.FindAllByUser(userId, page, tweetsPerPage)
	if err != nil {
		return nil, err
	}
	if len(tweets) == 0 {
		return nil, fmt.Errorf("no tweets found for user ID: %s", userId)
	}
	return tweets, nil
}
