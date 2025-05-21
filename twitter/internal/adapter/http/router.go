package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/twitter-tq/vinofsteel/internal/services"
)

// Router represents the HTTP router for the application
type Router struct {
	server *http.Server
}

// NewRouter creates a new Router with the given services
type NewRouterParams struct {
	UserService  services.UserServiceImpl
	TweetService services.TweetServiceImpl
	FollowService services.FollowServiceImpl
}

func NewRouter(params NewRouterParams) *Router {
	mux := http.NewServeMux()

	// User routes
	userHandlers := UserHandlers{params.UserService}
	mux.HandleFunc("POST /users", userHandlers.CreateUser)

	// Tweet routes
	tweetHandlers := TweetHandlers{params.TweetService}
	mux.HandleFunc("POST /tweets/{creator_id}", tweetHandlers.CreateTweet)
	mux.HandleFunc("GET /tweets/{creator_id}", tweetHandlers.GetTweetsByCreator)

	// Follow routes
	followHandlers := FollowHandlers{params.FollowService}
	mux.HandleFunc("POST /follows/{followed_id}", followHandlers.FollowUser)
	mux.HandleFunc("DELETE /follows/{followed_id}", followHandlers.UnfollowUser)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	return &Router{
		server: &server,
	}
}

// Run starts the HTTP server on the given address
func (r *Router) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting HTTP server", "address", r.server.Addr)
	return r.server.ListenAndServe()
}
