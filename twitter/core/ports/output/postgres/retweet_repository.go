package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
)

type RetweetRepository interface {
	HasRetweet(ctx context.Context, params RetweetHasParams) (bool, error)
	AddRetweet(ctx context.Context, params RetweetAddParams) (*models.Retweet, error)
}

// PostgresRetweetRepository implements RetweetRepository interface with PostgreSQL
type PostgresRetweetRepository struct {
	db *sql.DB
}

// NewPostgresRetweetRepository creates a new PostgreSQL retweet repository
func NewPostgresRetweetRepository(db *sql.DB) RetweetRepository {
	return &PostgresRetweetRepository{
		db: db,
	}
}

// HasRetweet checks if a user has retweeted a specific tweet
type RetweetHasParams struct {
	UserID  uuid.UUID `json:"user_id"`
	TweetID uuid.UUID `json:"tweet_id"`
}

func (r *PostgresRetweetRepository) HasRetweet(ctx context.Context, params RetweetHasParams) (bool, error) {
	slog.InfoContext(ctx, "Checking if user has retweeted tweet", "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
	query := `SELECT EXISTS(SELECT 1 FROM retweets WHERE user_id = $1 AND tweet_id = $2 AND deleted_at IS NULL);`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, params.UserID, params.TweetID).Scan(&exists)
	if err != nil {
		slog.ErrorContext(ctx, "Error checking if user has retweeted tweet", "error", err, "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
		return false, err
	}

	return exists, nil
}

// AddRetweet adds a new retweet to the database
type RetweetAddParams struct {
	UserID          uuid.UUID `json:"user_id"`
	TweetID         uuid.UUID `json:"tweet_id"`
}

func (r *PostgresRetweetRepository) AddRetweet(ctx context.Context, params RetweetAddParams) (*models.Retweet, error) {
	slog.InfoContext(ctx, "Adding retweet to database", "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
	query := `INSERT INTO retweets (user_id, tweet_id) VALUES ($1, $2) RETURNING id, created_at, updated_at, user_id, tweet_id;`

	var retweet models.Retweet
	err := r.db.QueryRowContext(
		ctx,
		query,
		params.UserID,
		params.TweetID,
	).Scan(
		&retweet.ID,
		&retweet.CreatedAt,
		&retweet.UpdatedAt,
		&retweet.UserID,
		&retweet.TweetID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error adding retweet to db", "error", err, "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
		return nil, err
	}

	return &retweet, nil
}