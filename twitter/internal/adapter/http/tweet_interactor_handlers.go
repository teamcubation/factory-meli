package http

import (
	"log/slog"
	"net/http"

	"github.com/twitter-tq/vinofsteel/internal/services"
)

type TweetInteractorHandlers struct {
	service services.TweetInteractorServices
}

// LikeTweet handles the HTTP request to like a tweet.
func (h TweetInteractorHandlers) LikeTweet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "Calling handler to like a tweet", "layer", "handler")

	err := h.service.Like(ctx, r)
	if err != nil {
		switch e := err.(type) {
		case services.ServiceError:
			respondWithError(ctx, w, err.(services.ServiceError).Code(), e.Error())
		default:
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(ctx, w, http.StatusOK, map[string]string{"message": "Tweet liked successfully"})
}

// UnlikeTweet handles the HTTP request to unlike a tweet.
func (h TweetInteractorHandlers) UnlikeTweet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "Calling handler to unlike a tweet", "layer", "handler")

	err := h.service.Unlike(ctx, r)
	if err != nil {
		switch e := err.(type) {
		case services.ServiceError:
			respondWithError(ctx, w, err.(services.ServiceError).Code(), e.Error())
		default:
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(ctx, w, http.StatusNoContent, "")
}

// RetweetTweet handles the HTTP request to retweet a tweet.
func (h TweetInteractorHandlers) RetweetTweet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "Calling handler to retweet a tweet", "layer", "handler")

	err := h.service.Retweet(ctx, r)
	if err != nil {
		switch e := err.(type) {
		case services.ServiceError:
			respondWithError(ctx, w, err.(services.ServiceError).Code(), e.Error())
		default:
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(ctx, w, http.StatusOK, map[string]string{"message": "tweet retweeted successfully"})
}