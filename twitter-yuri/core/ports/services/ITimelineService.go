package services

import "Yuri/twitter/core/models"

type ITimelineService interface {
	GetTimeline(userID string) ([]*models.Tweet, error)
}
