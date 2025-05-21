package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/internal/services"
)

// TweetHandlers handles HTTP requests for tweet operations
type TweetHandlers struct {
	service services.TweetServiceImpl
}

// Internal tweet type to return json correctly
type Tweet struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Post      string    `json:"post"`
	CreatorID uuid.UUID `json:"creator_id"`
}

func (h TweetHandlers) GetTweetsByCreator(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	creatorIDStr := r.PathValue("creator_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	tweets, err := h.service.GetAllUserTweets(ctx, creatorIDStr, limitStr, offsetStr)
	if err != nil {
		switch e := err.(type) {
		case services.ServiceError:
			respondWithError(ctx, w, err.(services.ServiceError).Code(), e.Error())
		default:
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(ctx, w, http.StatusOK, modelTweetsToTweets(tweets))
}

func (h TweetHandlers) CreateTweet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creatorIDStr := r.PathValue("creator_id")

	type parameters struct {
		Post string `json:"post"`
	}

	// Parse request body
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(ctx, w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create tweet using the service
	tweet, err := h.service.CreateTweet(ctx, params.Post, creatorIDStr)
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
	respondWithJSON(ctx, w, http.StatusCreated, modelTweetToTweet(*tweet))
}

// Utilities
func modelTweetToTweet(tweet models.Tweet) Tweet {
	return Tweet{
		ID:        tweet.ID,
		CreatedAt: tweet.CreatedAt,
		UpdatedAt: tweet.UpdatedAt,
		Post:      tweet.Post,
		CreatorID: tweet.CreatorID,
	}
}

func modelTweetsToTweets(tweets []*models.Tweet) []Tweet {
	result := make([]Tweet, len(tweets))
	for i, tweet := range tweets {
		result[i] = modelTweetToTweet(*tweet)
	}
	return result
}
