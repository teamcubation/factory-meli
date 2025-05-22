package services

import "Yuri/twitter/core/models"

type ITimelineService interface {
	GetTimeline(userID string, page, tweetsPerPage int) ([]*models.Tweet, error)
}
