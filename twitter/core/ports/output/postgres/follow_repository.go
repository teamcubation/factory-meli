package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
)

type FollowRepository interface {
	Save(ctx context.Context, params FollowSaveParams) (*models.Follow, error)
	Delete(ctx context.Context, params FollowDeleteParams) error
	FindByIds(ctx context.Context, params FollowFindByIdsParams) (*models.Follow, error)
}

// PostgresFollowRepository implements FollowRepository interface with PostgreSQL
type PostgresFollowRepository struct {
	db *sql.DB
}

func NewPostgresFollowRepository(db *sql.DB) FollowRepository {
	return &PostgresFollowRepository{
		db: db,
	}
}

type FollowSaveParams struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
}

func (r *PostgresFollowRepository) Save(ctx context.Context, params FollowSaveParams) (*models.Follow, error) {
	slog.InfoContext(ctx, "Saving a new tweet to the database", "follower_id", params.FollowerID, "followed_id", params.FollowedID, "layer", "database")

	query := `INSERT INTO follows (follower_id, followed_id) VALUES ($1, $2) RETURNING id, created_at, updated_at, follower_id, followed_id;`

	var follow models.Follow
	err := r.db.QueryRowContext(
		ctx,
		query,
		params.FollowerID,
		params.FollowedID,
	).Scan(
		&follow.ID,
		&follow.CreatedAt,
		&follow.UpdatedAt,
		&follow.FollowerID,
		&follow.FollowedID,
	)

	if err != nil {
		slog.ErrorContext(ctx, "Error saving follow to db", "error", err, "layer", "database")
		return nil, err
	}

	return &follow, nil
}

type FollowDeleteParams struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
}

func (r *PostgresFollowRepository) Delete(ctx context.Context, params FollowDeleteParams) error {
	slog.InfoContext(ctx, "Deleting a follow from the database", "follower_id", params.FollowerID, "followed_id", params.FollowedID, "layer", "database")
	query := `UPDATE follows SET deleted_at = NOW() WHERE follower_id = $1 AND followed_id = $2 AND deleted_at IS NULL;`

	result, err := r.db.ExecContext(ctx, query, params.FollowerID, params.FollowedID)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting follow", "error", err, "layer", "database")
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Error getting rows affected while deleting a follow", "error", err, "layer", "database")
		return err
	}

	if rows == 0 {
		slog.ErrorContext(ctx, "Follow relationship not found", "error", err, "layer", "database")
		return errors.New("follow relationship not found")
	}

	return nil
}

type FollowFindByIdsParams struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
}

func (r *PostgresFollowRepository) FindByIds(ctx context.Context, params FollowFindByIdsParams) (*models.Follow, error) {
	slog.InfoContext(ctx, "Getting a follow from the database", "follower_id", params.FollowerID, "followed_id", params.FollowedID, "layer", "database")
	query := `SELECT id, created_at, updated_at, follower_id, followed_id FROM follows WHERE follower_id = $1 AND followed_id = $2 AND deleted_at IS NULL;`

	var follow models.Follow
	err := r.db.QueryRowContext(ctx, query, params.FollowerID, params.FollowedID).Scan(
		&follow.ID,
		&follow.CreatedAt,
		&follow.UpdatedAt,
		&follow.FollowerID,
		&follow.FollowedID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error querying row while selecting a follow", "error", err, "layer", "database")
		return nil, err
	}

	return &follow, nil
}

type FollowFindFollowingParams struct {
	UserID uuid.UUID `json:"user_id"`
}
