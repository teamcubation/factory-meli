package service

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/out"
	"Yuri/twitter/core/ports/services"
)

type TimelineServiceImpl struct {
	repoTweet out.ITweetRepository
	repoUser  out.IUserRepository
}

func NewTimelineService(repoTweet out.ITweetRepository, repoUser out.IUserRepository) services.ITimelineService {
	return &TimelineServiceImpl{
		repoTweet: repoTweet,
		repoUser:  repoUser,
	}
}

func (s *TimelineServiceImpl) GetTimeline(userID string) ([]*models.Tweet, error) {
	panic("unimplemented")
}
