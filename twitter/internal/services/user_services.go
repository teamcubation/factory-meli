package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

// UserServices defines the interface for user-related interactions
type UserServices interface {
	CreateUser(ctx context.Context, r *http.Request) (*models.User, error)
	GetUserTimeline(ctx context.Context, r *http.Request) ([]*models.UserWithTweets, error)
}

type UserServiceImpl struct {
	repository postgres.UserRepository
}

func NewUserService(repo postgres.UserRepository) UserServices {
	return UserServiceImpl{
		repository: repo,
	}
}

func (s UserServiceImpl) CreateUser(ctx context.Context, r *http.Request) (*models.User, error) {
	slog.InfoContext(ctx, "Calling service to create a new user", "layer", "service")

	// Decoding body
	type parameters struct {
		Email string `json:"email"`
	}

	// Parse request body
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		slog.ErrorContext(ctx, "Invalid request payload for CreateUser", "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	// Validating email
	email := params.Email
	email = strings.TrimSpace(email)
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if email == "" {
		slog.ErrorContext(ctx, "Email is an obligatory field for CreateUser", "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "email is an obligatory field"}
	}

	if !emailRegex.MatchString(email) {
		slog.ErrorContext(ctx, "Invalid email format for CreateUser", "email", email, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid email format"}
	}

	// Verifying if the user email already has an account associated to it
	slog.InfoContext(ctx, "Verifying if user email already exists", "email", email, "layer", "service")
	user, err := s.repository.FindByEmail(ctx, postgres.UserFindByEmailParams{
		Email: email,
	})
	if err != nil && err != sql.ErrNoRows {
		slog.ErrorContext(ctx, "Error checking for existing user by email", "email", email, "error", err, "layer", "service")
		return nil, err
	}

	if user != nil {
		slog.ErrorContext(ctx, "User with this email already exists", "email", email, "layer", "service")
		return nil, ServiceError{http.StatusConflict, "there already exists an user with this email"}
	}

	// Saving new user's info
	slog.InfoContext(ctx, "Saving new user to the database", "email", email, "layer", "service")
	user, err = s.repository.Save(ctx, postgres.UserSaveParams{
		Email: email,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error saving new user", "email", email, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully created new user", "user_id", user.ID, "email", user.Email, "layer", "service")
	return user, nil
}

func (s UserServiceImpl) GetUserTimeline(ctx context.Context, r *http.Request) ([]*models.UserWithTweets, error) {
	slog.InfoContext(ctx, "Calling service to get user timeline", "layer", "service")

	userIDStr := r.PathValue("user_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	tweetNumStr := r.URL.Query().Get("tweet_num")

	// Validating parameters
	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
			slog.ErrorContext(ctx, "Invalid limit parameter for GetUserTimeline", "limit_str", limitStr, "error", err, "layer", "service")
			return nil, ServiceError{http.StatusBadRequest, "invalid limit"}
		}
	}

	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		} else {
			slog.ErrorContext(ctx, "Invalid offset parameter for GetUserTimeline", "offset_str", offsetStr, "error", err, "layer", "service")
			return nil, ServiceError{http.StatusBadRequest, "invalid offset"}
		}
	}

	tweetNum := 10
	if tweetNumStr != "" {
		if parsedTweetNum, err := strconv.Atoi(tweetNumStr); err == nil && parsedTweetNum >= 0 {
			tweetNum = parsedTweetNum
		} else {
			slog.ErrorContext(ctx, "Invalid tweet_num parameter for GetUserTimeline", "tweet_num_str", tweetNumStr, "error", err, "layer", "service")
			return nil, ServiceError{http.StatusBadRequest, "invalid tweet_num"}
		}
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid user_id in GetUserTimeline request", "user_id_str", userIDStr, "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid user_id"}
	}

	// Verifying if the user has an account associated to it
	slog.InfoContext(ctx, "Verifying user existence for timeline retrieval", "user_id", userID, "layer", "service")
	user, err := s.repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "User not found for timeline retrieval", "user_id", userID, "layer", "service")
			return nil, ServiceError{http.StatusNotFound, "user not found"}
		}
		slog.ErrorContext(ctx, "Error finding user by ID for timeline retrieval", "user_id", userID, "error", err, "layer", "service")
		return nil, err
	}

	// Getting user's timeline
	slog.InfoContext(ctx, "Retrieving user timeline", "user_id", user.ID, "limit", limit, "offset", offset, "tweets_per_followed", tweetNum, "layer", "service")
	timeline, err := s.repository.GetTimeline(ctx, postgres.TimelineParams{
		UserID:            user.ID,
		TweetsPerFollowed: tweetNum,
		Offset:            offset,
		Limit:             limit,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error retrieving user timeline", "user_id", user.ID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully retrieved user timeline", "user_id", user.ID, "num_users_with_tweets", len(timeline), "layer", "service")
	return timeline, nil
}
