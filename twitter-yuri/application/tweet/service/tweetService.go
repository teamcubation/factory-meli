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

	if err := s.repo.SaveTweet(tweet); err != nil {
		fmt.Println("Error saving tweet:", err)
		return err
	}
	return nil
}

func (s *TweetServiceImpl) ListTweetsByUserID(userId string) ([]*models.Tweet, error) {
	tweets, err := s.repo.FindAllByUser(userId)
	if err != nil {
		fmt.Println("Error fetching tweets:", err)
		return nil, err
	}

	if len(tweets) == 0 {
		return nil, fmt.Errorf("no tweets found for user ID: %s", userId)
	}
	return tweets, nil
}
