package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
)

// LikeRepository defines the interface for like persistence operations
type LikeRepository interface {
	HasLike(ctx context.Context, params LikeHasParams) (bool, error)
	AddLike(ctx context.Context, params LikeAddParams) (*models.Like, error)
	RemoveLike(ctx context.Context, params LikeRemoveParams) error
}

// PostgresLikeRepository implements LikeRepository interface with PostgreSQL
type PostgresLikeRepository struct {
	db *sql.DB
}

// NewPostgresLikeRepository creates a new PostgreSQL like repository
func NewPostgresLikeRepository(db *sql.DB) LikeRepository {
	return &PostgresLikeRepository{
		db: db,
	}
}

// HasLike checks if a user has liked a specific tweet
type LikeHasParams struct {
	UserID  uuid.UUID `json:"user_id"`
	TweetID uuid.UUID `json:"tweet_id"`
}

func (r *PostgresLikeRepository) HasLike(ctx context.Context, params LikeHasParams) (bool, error) {
	slog.InfoContext(ctx, "Checking if user has liked tweet", "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND tweet_id = $2 AND deleted_at IS NULL);`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, params.UserID, params.TweetID).Scan(&exists)
	if err != nil {
		slog.ErrorContext(ctx, "Error checking if user has liked tweet", "error", err, "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
		return false, err
	}

	return exists, nil
}

// AddLike adds a new like to the database
type LikeAddParams struct {
	UserID  uuid.UUID `json:"user_id"`
	TweetID uuid.UUID `json:"tweet_id"`
}

func (r *PostgresLikeRepository) AddLike(ctx context.Context, params LikeAddParams) (*models.Like, error) {
	slog.InfoContext(ctx, "Adding like to database", "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
	query := `INSERT INTO likes (user_id, tweet_id) VALUES ($1, $2) RETURNING id, created_at, updated_at, user_id, tweet_id;`

	var like models.Like
	err := r.db.QueryRowContext(
		ctx,
		query,
		params.UserID,
		params.TweetID,
	).Scan(
		&like.ID,
		&like.CreatedAt,
		&like.UpdatedAt,
		&like.UserID,
		&like.TweetID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error adding like to db", "error", err, "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
		return nil, err
	}

	return &like, nil
}

// RemoveLike soft deletes a like from the database
type LikeRemoveParams struct {
	UserID  uuid.UUID `json:"user_id"`
	TweetID uuid.UUID `json:"tweet_id"`
}

func (r *PostgresLikeRepository) RemoveLike(ctx context.Context, params LikeRemoveParams) error {
	slog.InfoContext(ctx, "Removing like from database", "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
	query := `UPDATE likes SET deleted_at = NOW(), updated_at = NOW() WHERE user_id = $1 AND tweet_id = $2 AND deleted_at IS NULL;`

	result, err := r.db.ExecContext(ctx, query, params.UserID, params.TweetID)
	if err != nil {
		slog.ErrorContext(ctx, "Error removing like from db", "error", err, "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Error getting rows affected when removing like", "error", err, "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
		return err
	}

	if rowsAffected == 0 {
		slog.WarnContext(ctx, "No like found to remove", "user_id", params.UserID, "tweet_id", params.TweetID, "layer", "database")
	}

	return nil
}
