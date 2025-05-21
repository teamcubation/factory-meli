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
	slog.InfoContext(ctx, "Saving a new tweet to the database", "follower_id", params.FollowerID, "followed_id", params.FollowedID)

	// // First check if the follow relationship already exists
	// exists, err := r.Exists(ctx, params.FollowerID, params.FollowedID)
	// if err != nil {
	// 	return nil, err
	// }

	// if exists {
	// 	return nil, errors.New("follow relationship already exists")
	// }

	// // Cannot follow yourself
	// if params.FollowerID == params.FollowedID {
	// 	return nil, errors.New("cannot follow yourself")
	// }

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
		slog.ErrorContext(ctx, "Error saving follow to db", "error", err)
		return nil, err
	}

	return &follow, nil
}

type FollowDeleteParams struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
}

func (r *PostgresFollowRepository) Delete(ctx context.Context, params FollowDeleteParams) error {
	query := `UPDATE follows SET deleted_at = NOW() WHERE follower_id = $1 AND followed_id = $2 AND deleted_at IS NULL;`

	result, err := r.db.ExecContext(ctx, query, params.FollowerID, params.FollowedID)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting follow", "error", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		slog.ErrorContext(ctx, "Error getting rows affected while deleting a follow", "error", err)
		return err
	}

	if rows == 0 {
		return errors.New("follow relationship not found")
	}

	return nil
}

type FollowFindFollowingParams struct {
	UserID uuid.UUID `json:"user_id"`
}

func (r *PostgresFollowRepository) FindFollowing(ctx context.Context, params FollowFindFollowingParams) ([]*models.User, error) {
	query := `
		SELECT u.id, u.created_at, u.updated_at, u.email FROM users u
			JOIN follows f ON u.id = f.followed_id
				WHERE f.follower_id = $1 AND f.deleted_at IS NULL AND u.deleted_at IS NULL;
	`

	rows, err := r.db.QueryContext(ctx, query, params.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "Error everyone that is following the user", "error", err, "user_id", params.UserID)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.Email,
		); err != nil {
			slog.ErrorContext(ctx, "Error scanning user", "error", err, "id", u.ID)
			return nil, err
		}
		users = append(users, &u)
	}

	if err := rows.Close(); err != nil {
		slog.ErrorContext(ctx, "Error closing rows", "error", err)
		return nil, err
	}

	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error when iterating rows", "error", err)
		return nil, err
	}

	return users, nil
}
