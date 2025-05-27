package ports

import "twitter-api-clone/core/models"

type UserRepository interface {
	Create(user *models.User) error
	Delete(userId string) error
}

type TweetRepository interface {
	Create(tweet *models.Tweet) error
	ListByUser(userId string) ([]models.Tweet, error)
	GetTimeline(userId string, limit int, offset int) ([]models.Timeline, error)
}

type FollowRepository interface {
	Follow(follow models.Follow) error
	Unfollow(follow models.Follow) error
}
