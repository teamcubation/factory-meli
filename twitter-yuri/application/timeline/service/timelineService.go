package service

import (
	"Yuri/twitter/core/models"
	"Yuri/twitter/core/ports/out"
	"Yuri/twitter/core/ports/services"
	"fmt"
	"sort"
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
	user, err := s.repoUser.GetById(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found with ID: %s", userID)
	}

	following := user.FollowUsers
	tweets := make([]*models.Tweet, 0)

	for _, userID := range following {
		userTweets, err := s.repoTweet.FindAllByUser(userID)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, userTweets...)
	}
	tweetsSorted := sortTweetsByDate(tweets)
	return tweetsSorted, nil
}

func sortTweetsByDate(tweets []*models.Tweet) []*models.Tweet {
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})
	return tweets
}
