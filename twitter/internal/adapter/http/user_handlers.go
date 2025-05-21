package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/internal/services"
)

// UserHandler handles HTTP requests for user operations
type UserHandlers struct {
	service services.UserServiceImpl
}

// Internal user types to return json correctly
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserWithTweets struct {
	User   `json:"user"`
	Tweets []*Tweet `json:"tweets"`
}

func (h UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type parameters struct {
		Email string `json:"email"`
	}

	// Parse request body
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(ctx, w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create user using the service
	user, err := h.service.CreateUser(ctx, params.Email)
	if err != nil {
		switch e := err.(type) {
		case services.ServiceError:
			respondWithError(ctx, w, err.(services.ServiceError).Code(), e.Error())
		default:
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Return success response
	respondWithJSON(ctx, w, http.StatusCreated, modelUserToUser(*user))
}

func (h UserHandlers) GetTimeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDStr := r.PathValue("user_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	tweetNumStr := r.URL.Query().Get("tweet_num")

	timeline, err := h.service.GetUserTimeline(ctx, userIDStr, limitStr, offsetStr, tweetNumStr)
	if err != nil {
		switch e := err.(type) {
		case services.ServiceError:
			respondWithError(ctx, w, err.(services.ServiceError).Code(), e.Error())
		default:
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(ctx, w, http.StatusOK, modelUsersWithTweetsToUsersWithTweets(timeline))
}

// Utilities
func modelUserToUser(user models.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}

// Converting the tweet slice values into pointers
func modelTweetsToTweetsFromValueSlice(tweets []models.Tweet) []*Tweet {
	result := make([]*Tweet, len(tweets))
	for i, tweet := range tweets {
		newTweet := modelTweetToTweet(tweet)
		result[i] = &newTweet
	}
	return result
}

func modelUserWithTweetsToUserWithTweets(user models.UserWithTweets) UserWithTweets {
	return UserWithTweets{
		User:   modelUserToUser(user.User),
		Tweets: modelTweetsToTweetsFromValueSlice(user.Tweets),
	}
}

func modelUsersWithTweetsToUsersWithTweets(users []*models.UserWithTweets) []UserWithTweets {
	result := make([]UserWithTweets, len(users))
	for i, user := range users {
		result[i] = modelUserWithTweetsToUserWithTweets(*user)
	}
	return result
}
