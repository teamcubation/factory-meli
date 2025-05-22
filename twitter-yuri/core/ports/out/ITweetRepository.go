package out

import "Yuri/twitter/core/models"

type ITweetRepository interface {
	SaveTweet(tweet *models.Tweet) error
	FindAllByUser(userId string, page, tweetsPerPage int) ([]*models.Tweet, error)
	FindAllByUserTimeline(userIds []string, page, tweetsPerPage int) ([]*models.Tweet, error)
}
