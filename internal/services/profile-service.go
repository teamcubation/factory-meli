package services

import (
	"tweet-tq-rafael/core/dtos"
	"tweet-tq-rafael/core/models"
	"tweet-tq-rafael/core/ports"
)

type UserService struct {
	Repo ports.UserRepository
}

func (s UserService) Create(createUserRequestDTO dtos.CreateUserRequestDTO) (models.User, error) {
	user, err := s.Repo.Create(createUserRequestDTO)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s UserService) Timeline(timeLineRequestDTO dtos.TimelineRequestDTO) ([]models.Tweet, error) {
	tweets, err := s.Repo.Timeline(timeLineRequestDTO)
	if err != nil {
		return []models.Tweet{}, err
	}
	return tweets, nil
}

func (s UserService) Follow(followRequestDTO dtos.FollowRequestDTO) error {
	err := s.Repo.Follow(followRequestDTO)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) Unfollow(unfollowRequestDTO dtos.UnfollowRequestDTO) error {

	err := s.Repo.Unfollow(unfollowRequestDTO)
	if err != nil {
		return err
	}

	return nil
}
