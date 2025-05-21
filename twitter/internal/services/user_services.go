package services

import (
	"context"
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

type UserServiceImpl struct {
	repository postgres.UserRepository
}

func NewUserService(repo postgres.UserRepository) UserServiceImpl {
	return UserServiceImpl{
		repository: repo,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, email string) (*models.User, error) {
	// Validating email
	email = strings.TrimSpace(email)
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if email == "" {
		return nil, ServiceError{http.StatusBadRequest, "email is an obligatory field"}
	}

	if !emailRegex.MatchString(email) {
		return nil, ServiceError{http.StatusBadRequest, "invalid email format"}
	}

	// Verifying if the user email already has an account associated to it
	user, err := s.repository.FindByEmail(ctx, postgres.UserFindByEmailParams{
		Email: email,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if user != nil {
		return nil, ServiceError{http.StatusBadRequest, "there already exists an user with this email"}
	}

	// Saving new user's info
	user, err = s.repository.Save(ctx, postgres.UserSaveParams{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetUserTimeline(ctx context.Context, userIDStr, limitStr, offsetStr, tweetNumStr string) ([]*models.UserWithTweets, error) {
	// Validating parameters
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

	tweetNum := 10
	if tweetNumStr != "" {
		if parsedTweetNum, err := strconv.Atoi(tweetNumStr); err == nil && parsedTweetNum >= 0 {
			tweetNum = parsedTweetNum
		} else {
			return nil, ServiceError{http.StatusBadRequest, "invalid tweet_num"}
		}
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
			return nil, ServiceError{http.StatusBadRequest, "invalid user_id"}
	}

	// Verifying if the user email has an account associated to it
	user, err := s.repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ServiceError{http.StatusNotFound, "user not found"}
		}
		return nil, err
	}

	// Getting user's timeline
	timeline, err := s.repository.GetTimeline(ctx, postgres.TimelineParams{
		UserID:            user.ID,
		TweetsPerFollowed: tweetNum,
		Offset:            offset,
		Limit:             limit,
	})
	if err != nil {
		return nil, err
	}

	return timeline, nil
}
