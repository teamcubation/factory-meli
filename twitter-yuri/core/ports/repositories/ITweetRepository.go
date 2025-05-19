package repositories

import "Yuri/twitter/core/models"

type ITweetRepository interface {
	SaveTweet(tweet *models.Tweet) error
	FindAllByUser(id int) ([]*models.Tweet, error)
}
