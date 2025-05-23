package ports

import (
	"tweet-tq-rafael/core/dtos"
	"tweet-tq-rafael/core/models"
)

type TweetRepository interface {
	Create(dtos.CreateTweetRequestDTO) (models.Tweet, error)
}

type UserRepository interface {
	Create(dtos.CreateUserRequestDTO) (models.User, error)
	Follow(dtos.FollowRequestDTO) error
	Unfollow(dtos.UnfollowRequestDTO) error
	Timeline(dtos.TimelineRequestDTO) ([]models.Tweet, error)
}

type Repositories struct {
	UserRepository  UserRepository
	TweetRepository TweetRepository
}
