package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog" // Import slog
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)


// TweetServices defines the interface for tweet-related interactions
type TweetServices interface {
	GetAllUserTweets(ctx context.Context, r *http.Request) ([]*models.Tweet, error)
	CreateTweet(ctx context.Context, r *http.Request) (*models.Tweet, error)
}

type TweetServiceImpl struct {
	t_repository postgres.TweetRepository
	u_repository postgres.UserRepository
}

func NewTweetService(tweetRepo postgres.TweetRepository, userRepo postgres.UserRepository) TweetServices {
	return TweetServiceImpl{
		t_repository: tweetRepo,
		u_repository: userRepo,
	}
}

func (s TweetServiceImpl) GetAllUserTweets(ctx context.Context, r *http.Request) ([]*models.Tweet, error) {
	slog.InfoContext(ctx, "Calling service to get all user tweets", "layer", "service")

	// Validating parameter
	limit := 10
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
			slog.ErrorContext(ctx, "Invalid limit parameter for GetAllUserTweets", "limit_str", limitStr, "error", err, "layer", "service")
			return nil, ServiceError{http.StatusBadRequest, "invalid limit"}
		}
	}

	offset := 0
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		} else {
			slog.ErrorContext(ctx, "Invalid offset parameter for GetAllUserTweets", "offset_str", offsetStr, "error", err, "layer", "service")
			return nil, ServiceError{http.StatusBadRequest, "invalid offset"}
		}
	}

	creatorIDStr := r.PathValue("creator_id")
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid creator_id in GetAllUserTweets request", "creator_id_str", creatorIDStr, "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid creator_id"}
	}

	slog.InfoContext(ctx, "Finding user by ID to get all tweets", "creator_id", creatorID, "layer", "service")
	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: creatorID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "User not found for GetAllUserTweets", "creator_id", creatorID, "layer", "service")
			return nil, ServiceError{http.StatusNotFound, "user not found"}
		}
		slog.ErrorContext(ctx, "Error finding user by ID for GetAllUserTweets", "creator_id", creatorID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Retrieving all tweets for user", "creator_id", creatorID, "limit", limit, "offset", offset, "layer", "service")
	tweets, err := s.t_repository.FindAllTweetsByUserId(ctx, postgres.TweetFindAllTweetsByUserId{
		CreatorID: creatorID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error retrieving all tweets by user ID", "creator_id", creatorID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully retrieved all user tweets", "creator_id", creatorID, "num_tweets", len(tweets), "layer", "service")
	return tweets, nil
}

func (s TweetServiceImpl) CreateTweet(ctx context.Context, r *http.Request) (*models.Tweet, error) {
	slog.InfoContext(ctx, "Calling service to create a new tweet", "layer", "service")

	// Parse request body
	type parameters struct {
		Post string `json:"post"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		slog.ErrorContext(ctx, "Invalid request payload for CreateTweet", "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	// Validating parameters
	post := strings.TrimSpace(params.Post)
	if post == "" {
		slog.ErrorContext(ctx, "Post content is empty in CreateTweet request", "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "post is an obligatory fields"}
	}

	if len(post) > 280 {
		slog.ErrorContext(ctx, "Post content exceeds character limit in CreateTweet request", "post_length", len(post), "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "post character limit is 280"}
	}

	creatorIDStr := r.PathValue("creator_id")
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid creator_id in CreateTweet request", "creator_id_str", creatorIDStr, "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid creator_id"}
	}

	slog.InfoContext(ctx, "Finding user by ID to create tweet", "creator_id", creatorID, "layer", "service")
	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: creatorID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "User not found for CreateTweet", "creator_id", creatorID, "layer", "service")
			return nil, ServiceError{http.StatusNotFound, "user not found"}
		}
		slog.ErrorContext(ctx, "Error finding user by ID for CreateTweet", "creator_id", creatorID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Saving new tweet to the database", "creator_id", creatorID, "post_length", len(post), "layer", "service")
	tweet, err := s.t_repository.Save(ctx, postgres.TweetSaveParams{
		Post:      post,
		CreatorID: creatorID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error saving new tweet", "creator_id", creatorID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully created new tweet", "tweet_id", tweet.ID, "creator_id", creatorID, "layer", "service")
	return tweet, nil
}
