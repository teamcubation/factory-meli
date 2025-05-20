package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/twitter-tq/vinofsteel/internal/ports/output/postgres"
)

func main() {
	// Creating application context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Initializing environment variables
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load("../../.env"); err != nil {
			slog.ErrorContext(ctx, "Error loading .env file", "error", err)
			os.Exit(1)
		}
	}

	// Setting up db
	dbProvider := postgres.NewPostgresDatabaseProvider()
	defer dbProvider.Close()
	_, err := dbProvider.GetConnection(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting SQL connection", "error", err)
		os.Exit(1)
	}

	log.Println("everything alright")

	// Creating model repos
	// userRepo := postgres.NewPostgresUserRepository(db)
}