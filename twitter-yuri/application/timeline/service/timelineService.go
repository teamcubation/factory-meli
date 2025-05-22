package service

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/out"
	"Yuri/twitter/core/ports/services"
	"fmt"
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

func (s *TimelineServiceImpl) GetTimeline(userID string, page, tweetsPerPage int) ([]*models.Tweet, error) {
	user, err := s.repoUser.GetById(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found with ID: %s", userID)
	}
	userTweets, err := s.repoTweet.FindAllByUserTimeline(user.FollowUsers, page, tweetsPerPage)
	if err != nil {
		return nil, err
	}
	return userTweets, nil
}
