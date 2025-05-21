package services

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

type TweetServiceImpl struct {
	repository postgres.TweetRepository
}

func NewTweetService(repo postgres.TweetRepository) TweetServiceImpl {
	return TweetServiceImpl{
		repository: repo,
	}
}

func (s *TweetServiceImpl) CreateTweet(ctx context.Context, post, creatorIDStr string) (*models.Tweet, error) {
	// Validating parameters
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		return nil, ServiceError{http.StatusBadRequest, "invalid creator_id"}
	}

	post = strings.TrimSpace(post)

	if post == "" {
		return nil, ServiceError{http.StatusBadRequest, "post is an obligatory fields"}
	}

	if len(post) > 280 {
		return nil, ServiceError{http.StatusBadRequest, "post character limit is 280"}
	}

	tweet, err := s.repository.Save(ctx, postgres.TweetSaveParams{
		Post:      post,
		CreatorID: creatorID,
	})
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (s *TweetServiceImpl) GetAllUserTweets(ctx context.Context, creatorIDStr, limitStr, offsetStr string) ([]*models.Tweet, error) {
	// Validating parameters
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		return nil, ServiceError{http.StatusBadRequest, "invalid creator_id"}
	}

	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
		return nil, ServiceError{http.StatusBadRequest, "invalid limit"}
		}
	}

	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		} else {
		return nil, ServiceError{http.StatusBadRequest, "invalid offset"}
		}
	}

	tweet, err := s.repository.FindAllTweetsByUserId(ctx, postgres.TweetFindAllTweetsByUserId{
		CreatorID: creatorID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}

	return tweet, nil
}
