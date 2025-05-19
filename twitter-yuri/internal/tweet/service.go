package tweet

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/repositories"
	"Yuri/twitter/core/ports/services"
)

type ServiceImpl struct {
	repo   repositories.ITweetRepository
	nextID int
}

func NewTweetService(repo repositories.ITweetRepository) services.ITweetService {
	return &ServiceImpl{
		repo:   repo,
		nextID: 1,
	}
}

func (s *ServiceImpl) CreateTweet(tweet *models.Tweet) error {
	panic("unimplemented")
}

func (s *ServiceImpl) ListTweetsByUserID(id int) ([]*models.Tweet, error) {
	panic("unimplemented")
}
