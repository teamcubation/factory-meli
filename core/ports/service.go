package ports

import (
	"tweet-tq-rafael/core/dtos"
	"tweet-tq-rafael/core/models"
)

type Services struct {
	TweetService TweetService
	UserService  UserService
}

type TweetService interface {
	Create(dtos.CreateTweetRequestDTO) (models.Tweet, error)
}

type UserService interface {
	Create(dtos.CreateUserRequestDTO) (models.User, error)
	Follow(dtos.FollowRequestDTO) error
	Unfollow(dtos.UnfollowRequestDTO) error
	Timeline(dtos.TimelineRequestDTO) ([]models.Tweet, error)
}
