package service

import (
	"time"
	"twitter-api-clone/core/models"
	"twitter-api-clone/core/ports"
	errors "twitter-api-clone/internal/shared"

	"github.com/google/uuid"
)

type TweetServiceImpl struct {
	repo ports.TweetRepository
}

func NewTweetService(repo ports.TweetRepository) ports.TweetService {
	return &TweetServiceImpl{
		repo: repo,
	}
}

func (u *TweetServiceImpl) Create(tweet *models.Tweet) error {
	if len(tweet.Content) > 280 {
		return errors.ErrContentTooLong
	}

	tweet.Id = uuid.NewString()
	tweet.CreatedAt = time.Now()

	if err := u.repo.Create(tweet); err != nil {
		return err
	}

	return nil
}

func (u *TweetServiceImpl) ListByUser(userId string) ([]models.Tweet, error) {
	return u.repo.ListByUser(userId)
}

func (u *TweetServiceImpl) GetTimeline(userId string, limit int, page int) ([]models.Timeline, error) {
	offset := (page - 1) * limit

	return u.repo.GetTimeline(userId, limit, offset)
}
