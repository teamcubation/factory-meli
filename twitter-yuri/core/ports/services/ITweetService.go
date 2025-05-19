package services

import "Yuri/twitter/core/models"

type ITweetService interface {
	CreateTweet(tweet *models.Tweet) error
	ListTweetsByUserID(id int) ([]*models.Tweet, error)
}
