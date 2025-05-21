package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
)

// TweetRepository defines the interface for tweet persistence operations
type TweetRepository interface {
	Save(ctx context.Context, params TweetSaveParams) (*models.Tweet, error)
	FindByID(ctx context.Context, params TweetFindByIDParams) (*models.Tweet, error)
	FindAllTweetsByUserId(ctx context.Context, params TweetFindAllTweetsByUserId) ([]*models.Tweet, error)
}

// PostgresTweetRepository implements TweetRepository interface with PostgreSQL
type PostgresTweetRepository struct {
	db *sql.DB
}

// NewPostgresTweetRepository creates a new PostgreSQL user repository
func NewPostgresTweetRepository(db *sql.DB) TweetRepository {
	return &PostgresTweetRepository{
		db: db,
	}
}

// Save inserts a new tweet into the database
type TweetSaveParams struct {
	Post      string    `json:"post"`
	CreatorID uuid.UUID `json:"creator_id"`
}

func (r *PostgresTweetRepository) Save(ctx context.Context, params TweetSaveParams) (*models.Tweet, error) {
	slog.InfoContext(ctx, "Saving a new tweet to the database", "post", params.Post, "creator_id", params.CreatorID)
	query := `INSERT INTO tweets (post, creator_id) VALUES ($1, $2) RETURNING id, created_at, updated_at, post, creator_id;`

	var tweet models.Tweet
	err := r.db.QueryRowContext(
		ctx,
		query,
		params.Post,
		params.CreatorID,
	).Scan(
		&tweet.ID,
		&tweet.CreatedAt,
		&tweet.UpdatedAt,
		&tweet.Post,
		&tweet.CreatorID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error saving tweet to db", "error", err)
		return nil, err
	}

	return &tweet, nil
}

// FindById retrieves a user by their id
type TweetFindByIDParams struct {
	ID uuid.UUID `json:"id"`
}

func (r *PostgresTweetRepository) FindByID(ctx context.Context, params TweetFindByIDParams) (*models.Tweet, error) {
	slog.InfoContext(ctx, "Finding a tweet by id", "id", params.ID)
	query := `SELECT id, created_at, updated_at, post, creator_id FROM tweets WHERE id = $1 AND deleted_at IS NULL;`

	var tweet models.Tweet
	err := r.db.QueryRowContext(ctx, query, params.ID).Scan(
		&tweet.ID,
		&tweet.CreatedAt,
		&tweet.UpdatedAt,
		&tweet.Post,
		&tweet.CreatorID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error saving getting tweet from db by ID", "error", err, "id", params.ID)
		return nil, err
	}

	return &tweet, nil
}

type TweetFindAllTweetsByUserId struct {
	CreatorID uuid.UUID `json:"creator_id"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

func (r *PostgresTweetRepository) FindAllTweetsByUserId(ctx context.Context, params TweetFindAllTweetsByUserId) ([]*models.Tweet, error) {
	slog.InfoContext(ctx, "Finding all tweets of a user by their id", "creator_id", params.CreatorID)

	query := `
		SELECT id, created_at, updated_at, post, creator_id
			FROM tweets WHERE creator_id = $1
				ORDER BY created_at DESC LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.QueryContext(ctx, query, params.CreatorID, params.Limit, params.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error querying db to find all tweets of a user", "error", err, "creator_id", params.CreatorID)
		return nil, err
	}
	defer rows.Close()

	var tweets []*models.Tweet
	for rows.Next() {
		var t models.Tweet
		if err := rows.Scan(
			&t.ID,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.Post,
			&t.CreatorID,
		); err != nil {
			slog.ErrorContext(ctx, "Error scanning tweet", "error", err, "id", t.ID)
			return nil, err
		}
		tweets = append(tweets, &t)
	}

	if err := rows.Close(); err != nil {
		slog.ErrorContext(ctx, "Error closing rows", "error", err)
		return nil, err
	}

	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error when iterating rows", "error", err)
		return nil, err
	}

	return tweets, nil
}
