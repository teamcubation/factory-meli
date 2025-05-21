package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

type TweetServiceImpl struct {
	t_repository postgres.TweetRepository
	u_repository postgres.UserRepository
}

func NewTweetService(tweetRepo postgres.TweetRepository, userRepo postgres.UserRepository) TweetServiceImpl {
	return TweetServiceImpl{
		t_repository: tweetRepo,
		u_repository: userRepo,
	}
}

func (s *TweetServiceImpl) GetAllUserTweets(ctx context.Context, r *http.Request) ([]*models.Tweet, error) {	
	// Validating parameters
	creatorIDStr := r.PathValue("creator_id")
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		return nil, ServiceError{http.StatusBadRequest, "invalid creator_id"}
	}

	limit := 10
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
		return nil, ServiceError{http.StatusBadRequest, "invalid limit"}
		}
	}

	offset := 0
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		} else {
		return nil, ServiceError{http.StatusBadRequest, "invalid offset"}
		}
	}

	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: creatorID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ServiceError{http.StatusNotFound, "user not found"}
		}
		return nil, err
	}

	tweet, err := s.t_repository.FindAllTweetsByUserId(ctx, postgres.TweetFindAllTweetsByUserId{
		CreatorID: creatorID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (s *TweetServiceImpl) CreateTweet(ctx context.Context, r *http.Request) (*models.Tweet, error) {
	// Parse request body
	type parameters struct {
		Post string `json:"post"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	// Validating parameters
	post := strings.TrimSpace(params.Post)
	if post == "" {
		return nil, ServiceError{http.StatusBadRequest, "post is an obligatory fields"}
	}

	if len(post) > 280 {
		return nil, ServiceError{http.StatusBadRequest, "post character limit is 280"}
	}

	creatorIDStr := r.PathValue("creator_id")
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		return nil, ServiceError{http.StatusBadRequest, "invalid creator_id"}
	}

	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: creatorID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ServiceError{http.StatusNotFound, "user not found"}
		}
		return nil, err
	}

	tweet, err := s.t_repository.Save(ctx, postgres.TweetSaveParams{
		Post:      post,
		CreatorID: creatorID,
	})
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

