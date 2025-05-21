package out

import "Yuri/twitter/core/models"

type ITweetRepository interface {
	SaveTweet(tweet *models.Tweet) error
	FindAllByUser(userId string) ([]*models.Tweet, error)
}
