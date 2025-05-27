package ports

import "twitter-api-clone/core/models"

type TweetService interface {
	Create(tweet *models.Tweet) error
	ListByUser(userId string) ([]models.Tweet, error)
	GetTimeline(userId string, limit int, page int) ([]models.Timeline, error)
}

type FollowService interface {
	Follow(follow models.Follow) error
	Unfollow(follow models.Follow) error
}

type UserService interface {
	Create(user *models.User) error
	Delete(userId string) error
}
