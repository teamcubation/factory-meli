package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
	"github.com/twitter-tq/vinofsteel/internal/adapter/http"
	"github.com/twitter-tq/vinofsteel/internal/services"
	"github.com/twitter-tq/vinofsteel/pkg/logging"
)

func main() {
	// Creating application context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Initializing environment variables
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			slog.ErrorContext(ctx, "Error loading .env file", "error", err)
			os.Exit(1)
		}
	}

	// Setting up logging
	logging.SetupLogger(ctx)

	// Setting up db
	dbProvider := postgres.NewPostgresDatabaseProvider()
	defer dbProvider.Close()
	db, err := dbProvider.GetConnection(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting SQL connection", "error", err)
		os.Exit(1)
	}

	// Creating model services
	userRepo := postgres.NewPostgresUserRepository(db)
	userService := services.NewUserService(userRepo)

	tweetRepo := postgres.NewPostgresTweetRepository(db)
	tweetService := services.NewTweetService(tweetRepo, userRepo)

	followRepo := postgres.NewPostgresFollowRepository(db)
	followService := services.NewFollowService(followRepo, userRepo)

	likeRepo := postgres.NewPostgresLikeRepository(db)
	retweetRepo := postgres.NewPostgresRetweetRepository(db)
	tweetInteractorService := services.NewTweetInteractor(userRepo, likeRepo, tweetRepo, retweetRepo)

	router := http.NewRouter(http.NewRouterParams{
		UserService:             userService,
		TweetService:            tweetService,
		FollowService:           followService,
		TweetInteractorServices: tweetInteractorService,
	})
	router.Run(ctx)
}
