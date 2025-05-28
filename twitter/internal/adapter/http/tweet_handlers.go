package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/internal/services"
)

// TweetHandlers handles HTTP requests for tweet operations
type TweetHandlers struct {
	service services.TweetServices
}

// Internal tweet type to return json correctly
type Tweet struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Post      string     `json:"post"`
	CreatorID uuid.UUID  `json:"creator_id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
}

func (h TweetHandlers) GetTweetsByCreator(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "Calling handler to get tweets of an user", "layer", "handler")

	tweets, err := h.service.GetAllUserTweets(ctx, r)
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
	slog.InfoContext(ctx, "Calling handler to create a new tweet", "layer", "handler")

	// Create tweet using the service
	tweet, err := h.service.CreateTweet(ctx, r)
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

func (h TweetHandlers) CreateReplyTweet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "Calling handler to create a new reply to an existing tweet", "layer", "handler")

	// Create reply using the service
	tweet, err := h.service.ReplyToTweet(ctx, r)
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

func (h TweetHandlers) GetTweetThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "Calling handler to get thread based on a root tweet", "layer", "handler")

	// Getting flat tweets thread
	tweets, err := h.service.GetTweetThread(ctx, r)
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
	respondWithJSON(ctx, w, http.StatusCreated, modelTweetsToTweets(tweets))
}

// Utilities
func modelTweetToTweet(tweet models.Tweet) Tweet {
	var parentID *uuid.UUID
	if tweet.ParentID.Valid {
		parentID = &tweet.ParentID.UUID
	}

	return Tweet{
		ID:        tweet.ID,
		CreatedAt: tweet.CreatedAt,
		UpdatedAt: tweet.UpdatedAt,
		Post:      tweet.Post,
		CreatorID: tweet.CreatorID,
		ParentID:  parentID,
	}
}

func modelTweetsToTweets(tweets []*models.Tweet) []Tweet {
	result := make([]Tweet, len(tweets))
	for i, tweet := range tweets {
		result[i] = modelTweetToTweet(*tweet)
	}
	return result
}
