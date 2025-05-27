package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
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
	ReplyToTweet(ctx context.Context, r *http.Request) (*models.Tweet, error)
	GetTweetThread(ctx context.Context, r *http.Request) ([]*models.Tweet, error)
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

// This service returns all tweets of an user
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

// This service creates a new tweet
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

// This service handles the creation of a reply to an existing tweet.
func (s TweetServiceImpl) ReplyToTweet(ctx context.Context, r *http.Request) (*models.Tweet, error) {
	slog.InfoContext(ctx, "Calling service to reply to a tweet", "layer", "service")

	// Parse request body
	type parameters struct {
		Post string `json:"post"`
		CreatorID string `json:"creator_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		slog.ErrorContext(ctx, "Invalid request payload for ReplyToTweet", "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	// Validating parameters
	post := strings.TrimSpace(params.Post)
	if post == "" {
		slog.ErrorContext(ctx, "Post content is empty in ReplyToTweet request", "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "post is an obligatory field"}
	}

	if len(post) > 280 {
		slog.ErrorContext(ctx, "Post content exceeds character limit in ReplyToTweet request", "post_length", len(post), "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "post character limit is 280"}
	}

	// Extract parent tweet ID from path
	parentTweetIDStr := r.PathValue("tweet_id")
	parentTweetID, err := uuid.Parse(parentTweetIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid parent_tweet_id in ReplyToTweet request", "parent_tweet_id_str", parentTweetIDStr, "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid parent tweet ID"}
	}

	// Extract creator ID from path
	creatorIDStr := params.CreatorID
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid creator_id in ReplyToTweet request", "creator_id_str", creatorIDStr, "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid creator_id"}
	}

	slog.InfoContext(ctx, "Confirming existence of parent tweet", "parent_tweet_id", parentTweetID, "layer", "service")
	_, err = s.t_repository.FindByID(ctx, postgres.TweetFindByIDParams{
		ID: parentTweetID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Parent tweet not found for ReplyToTweet", "parent_tweet_id", parentTweetID, "layer", "service")
			return nil, ServiceError{http.StatusNotFound, "parent tweet not found"}
		}
		slog.ErrorContext(ctx, "Error finding parent tweet by ID for ReplyToTweet", "parent_tweet_id", parentTweetID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Finding user by ID to create reply tweet", "creator_id", creatorID, "layer", "service")
	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: creatorID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "User not found for ReplyToTweet", "creator_id", creatorID, "layer", "service")
			return nil, ServiceError{http.StatusNotFound, "user not found"}
		}
		slog.ErrorContext(ctx, "Error finding user by ID for ReplyToTweet", "creator_id", creatorID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Saving new reply tweet to the database", "creator_id", creatorID, "parent_tweet_id", parentTweetID, "post_length", len(post), "layer", "service")
	tweet, err := s.t_repository.SaveReply(ctx, postgres.TweetSaveReplyParams{
		Post:      post,
		CreatorID: creatorID,
		ParentID:  parentTweetID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error saving new reply tweet", "creator_id", creatorID, "parent_tweet_id", parentTweetID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully created new reply tweet", "tweet_id", tweet.ID, "creator_id", creatorID, "parent_tweet_id", parentTweetID, "layer", "service")
	return tweet, nil
}


// This service retrieves a complete thread of tweets, starting from a given root tweet ID.
func (s TweetServiceImpl) GetTweetThread(ctx context.Context, r *http.Request) ([]*models.Tweet, error) {
	slog.InfoContext(ctx, "Calling service to get a tweet thread", "layer", "service")

	rootTweetIDStr := r.PathValue("tweet_id")
	rootTweetID, err := uuid.Parse(rootTweetIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid root_tweet_id in GetTweetThread request", "root_tweet_id_str", rootTweetIDStr, "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid root tweet ID"}
	}

	slog.InfoContext(ctx, "Finding root tweet by ID to get thread", "root_tweet_id", rootTweetID, "layer", "service")
	_, err = s.t_repository.FindByID(ctx, postgres.TweetFindByIDParams{
		ID: rootTweetID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Root tweet not found for GetTweetThread", "root_tweet_id", rootTweetID, "layer", "service")
			return nil, ServiceError{http.StatusNotFound, "root tweet not found"}
		}
		slog.ErrorContext(ctx, "Error finding root tweet by ID for GetTweetThread", "root_tweet_id", rootTweetID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Fetching all replies for the root tweet", "root_tweet_id", rootTweetID, "layer", "service")
	// The repository returns a flat list of all tweets in the thread, including the root
	threadTweets, err := s.t_repository.FetchReplies(ctx, postgres.TweetFetchRepliesParams{
		RootTweetID: rootTweetID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error fetching replies for tweet thread", "root_tweet_id", rootTweetID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully retrieved tweet thread", "root_tweet_id", rootTweetID, "num_tweets_in_thread", len(threadTweets), "layer", "service")
	return threadTweets, nil
}