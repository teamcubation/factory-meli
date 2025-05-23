package services

import (
	"tweet-tq-rafael/core/dtos"
	"tweet-tq-rafael/core/models"
	"tweet-tq-rafael/core/ports"
)

type TweetService struct {
	Repo ports.TweetRepository
}

func (s TweetService) Create(createTweetRequestDTO dtos.CreateTweetRequestDTO) (models.Tweet, error) {
	tweet, err := s.Repo.Create(createTweetRequestDTO)
	if err != nil {
		return models.Tweet{}, err
	}
	return tweet, nil
}
