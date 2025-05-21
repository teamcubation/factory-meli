package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	Save(ctx context.Context, params UserSaveParams) (*models.User, error)
	FindByEmail(ctx context.Context, params UserFindByEmailParams) (*models.User, error)
	FindByID(ctx context.Context, params UserFindByIDParams) (*models.User, error)
	GetTimeline(ctx context.Context, params TimelineParams) ([]*models.UserWithTweets, error)
}

// PostgresUserRepository implements UserRepository interface with PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Save inserts a new user into the database
type UserSaveParams struct {
	Email string `json:"email"`
}

func (r *PostgresUserRepository) Save(ctx context.Context, params UserSaveParams) (*models.User, error) {
	slog.InfoContext(ctx, "Saving a new user to the database", "email", params.Email)
	query := `INSERT INTO users (email) VALUES ($1) RETURNING id, created_at, updated_at, email;`

	var user models.User
	err := r.db.QueryRowContext(
		ctx,
		query,
		params.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error saving user to db", "error", err)
		return nil, err
	}

	return &user, nil
}

// FindByEmail retrieves a user by their email
type UserFindByEmailParams struct {
	Email string `json:"email"`
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, params UserFindByEmailParams) (*models.User, error) {
	slog.InfoContext(ctx, "Finding an user by email", "email", params.Email)
	query := `SELECT id, created_at, updated_at, email FROM users WHERE email = $1 AND deleted_at IS NULL;`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, params.Email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting user from db by email", "error", err, "email", params.Email)
		return nil, err
	}

	return &user, nil
}

// FindByID retrieves a user by their ID
type UserFindByIDParams struct {
	ID uuid.UUID `json:"id"`
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, params UserFindByIDParams) (*models.User, error) {
	slog.InfoContext(ctx, "Finding an user by ID", "id", params.ID)
	query := `SELECT id, created_at, updated_at, email FROM users WHERE id = $1 AND deleted_at IS NULL;`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, params.ID).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting user from db by ID", "error", err, "id", params.ID)
		return nil, err
	}

	return &user, nil
}

// TimelineParams defines the parameters for retrieving a user's timeline
type TimelineParams struct {
	UserID            uuid.UUID `json:"user_id"`
	TweetsPerFollowed int       `json:"tweets_per_followed"`
	Limit             int       `json:"limit"`
	Offset            int       `json:"offset"`
}

func (r *PostgresUserRepository) GetTimeline(ctx context.Context, params TimelineParams) ([]*models.UserWithTweets, error) {
	slog.InfoContext(ctx, "Fetching timeline for user", "user_id", params.UserID, "tweets_per_followed", params.TweetsPerFollowed, "limit", params.Limit, "offset", params.Offset)

	// This query gets all users that the specified user follows and their tweets based on the business logic
	query := `
		WITH followed_users AS (
			SELECT followed_id 
			FROM follows 
			WHERE follower_id = $1 AND deleted_at IS NULL
		),
		tweets_from_followed AS (
			SELECT 
				u.id AS user_id, 
				u.email, 
				u.created_at AS user_created_at, 
				u.updated_at AS user_updated_at,
				t.id AS tweet_id, 
				t.created_at AS tweet_created_at, 
				t.updated_at AS tweet_updated_at, 
				t.post, 
				t.creator_id,
				-- Create a row number partitioned by creator to limit tweets per user
				ROW_NUMBER() OVER (PARTITION BY t.creator_id ORDER BY t.created_at DESC) AS tweet_rank
			FROM 
				users u
				JOIN tweets t ON u.id = t.creator_id
			WHERE 
				u.id IN (SELECT followed_id FROM followed_users)
				AND t.deleted_at IS NULL
				AND u.deleted_at IS NULL
		)
		SELECT 
			user_id, 
			email, 
			user_created_at, 
			user_updated_at,
			tweet_id, 
			tweet_created_at, 
			tweet_updated_at, 
			post, 
			creator_id
		FROM 
			tweets_from_followed
		WHERE 
			tweet_rank <= $2
		ORDER BY 
			tweet_created_at DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.db.QueryContext(ctx, query, params.UserID, params.TweetsPerFollowed, params.Limit, params.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error querying timeline", "error", err, "user_id", params.UserID)
		return nil, err
	}
	defer rows.Close()

	// Map to store users and their tweets
	userMap := make(map[uuid.UUID]*models.UserWithTweets)

	for rows.Next() {
		var userID uuid.UUID
		var email string
		var userCreatedAt, userUpdatedAt sql.NullTime
		var tweetID uuid.UUID
		var tweetCreatedAt, tweetUpdatedAt sql.NullTime
		var post string
		var creatorID uuid.UUID

		if err := rows.Scan(
			&userID,
			&email,
			&userCreatedAt,
			&userUpdatedAt,
			&tweetID,
			&tweetCreatedAt,
			&tweetUpdatedAt,
			&post,
			&creatorID,
		); err != nil {
			slog.ErrorContext(ctx, "Error scanning timeline row", "error", err)
			return nil, err
		}

		// Get or create user entry in the map
		userWithTweet, exists := userMap[userID]
		if !exists {
			userWithTweet = &models.UserWithTweets{
				User: models.User{
					ID:        userID,
					CreatedAt: userCreatedAt.Time,
					UpdatedAt: userUpdatedAt.Time,
					Email:     email,
				},
				Tweets: []models.Tweet{},
			}
			userMap[userID] = userWithTweet
		}

		// Add the tweet to the user's tweets
		tweet := models.Tweet{
			ID:        tweetID,
			CreatedAt: tweetCreatedAt.Time,
			UpdatedAt: tweetUpdatedAt.Time,
			Post:      post,
			CreatorID: creatorID,
		}
		userWithTweet.Tweets = append(userWithTweet.Tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error when iterating rows", "error", err)
		return nil, err
	}

	// Convert map to slice
	result := make([]*models.UserWithTweets, 0, len(userMap))
	for _, userWithTweets := range userMap {
		result = append(result, userWithTweets)
	}

	return result, nil
}
