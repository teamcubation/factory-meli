package service

import (
	"time"
	"twitter-api-clone/core/models"
	"twitter-api-clone/core/ports"
)

type FollowServiceImpl struct {
	repo ports.FollowRepository
}

func NewFollowService(repo ports.FollowRepository) ports.FollowService {
	return &FollowServiceImpl{
		repo: repo,
	}
}

func (f *FollowServiceImpl) Follow(follow models.Follow) error {
	follow.CreatedAt = time.Now()

	// TODO: add validation to prevent follow yourself
	// TODO: add validation to prevent follow the same user twice

	if err := f.repo.Follow(follow); err != nil {
		return err
	}

	return nil
}

func (f *FollowServiceImpl) Unfollow(follow models.Follow) error {
	if err := f.repo.Unfollow(follow); err != nil {
		return err
	}

	return nil
}
