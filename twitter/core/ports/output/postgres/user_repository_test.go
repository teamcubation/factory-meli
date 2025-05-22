package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// User Repository Tests
func TestUserSave_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	email := "test@example.com"
	userID := uuid.New()
	now := time.Now()

	// Set up expectations
	mock.ExpectQuery(`INSERT INTO users \(email\) VALUES \(\$1\) RETURNING id, created_at, updated_at, email`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "email"}).
			AddRow(userID, now, now, email))

	// Execute the method
	ctx := context.Background()
	params := UserSaveParams{Email: email}
	user, err := repo.Save(ctx, params)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSave_Error(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	email := "test@example.com"

	// Set up expectation for database error
	mock.ExpectQuery(`INSERT INTO users \(email\) VALUES \(\$1\) RETURNING id, created_at, updated_at, email`).
		WithArgs(email).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := UserSaveParams{Email: email}
	user, err := repo.Save(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSave_DuplicateEmailError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	email := "existing@example.com"

	// Set up expectation for a generic database error that contains the unique constraint message
	mock.ExpectQuery(`INSERT INTO users \(email\) VALUES \(\$1\) RETURNING id, created_at, updated_at, email`).
		WithArgs(email).
		WillReturnError(
			errors.New(`pq: duplicate key value violates unique constraint "users_email_key"`),
		)

	ctx := context.Background()
	params := UserSaveParams{Email: email}
	user, err := repo.Save(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, user)

	// Assert against the specific error if possible.
	assert.Equal(t, errors.New(`pq: duplicate key value violates unique constraint "users_email_key"`), err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFindByEmail_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	email := "test@example.com"
	userID := uuid.New()
	now := time.Now()

	// Set up expectations
	mock.ExpectQuery(`SELECT id, created_at, updated_at, email FROM users WHERE email = \$1 AND deleted_at IS NULL`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "email"}).
			AddRow(userID, now, now, email))

	// Execute the method
	ctx := context.Background()
	params := UserFindByEmailParams{Email: email}
	user, err := repo.FindByEmail(ctx, params)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFindByEmail_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	email := "notfound@example.com"

	// Set up expectation for no rows found
	mock.ExpectQuery(`SELECT id, created_at, updated_at, email FROM users WHERE email = \$1 AND deleted_at IS NULL`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	params := UserFindByEmailParams{Email: email}
	user, err := repo.FindByEmail(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFindByEmail_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	email := "test@example.com"

	// Set up expectation for database connection error
	mock.ExpectQuery(`SELECT id, created_at, updated_at, email FROM users WHERE email = \$1 AND deleted_at IS NULL`).
		WithArgs(email).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := UserFindByEmailParams{Email: email}
	user, err := repo.FindByEmail(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFindByID_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	userID := uuid.New()
	email := "test@example.com"
	now := time.Now()

	// Set up expectations
	mock.ExpectQuery(`SELECT id, created_at, updated_at, email FROM users WHERE id = \$1 AND deleted_at IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "email"}).
			AddRow(userID, now, now, email))

	// Execute the method
	ctx := context.Background()
	params := UserFindByIDParams{ID: userID}
	user, err := repo.FindByID(ctx, params)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFindByID_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	userID := uuid.New()

	// Set up expectation for no rows found
	mock.ExpectQuery(`SELECT id, created_at, updated_at, email FROM users WHERE id = \$1 AND deleted_at IS NULL`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	params := UserFindByIDParams{ID: userID}
	user, err := repo.FindByID(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFindByID_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	userID := uuid.New()

	// Set up expectation for database connection error
	mock.ExpectQuery(`SELECT id, created_at, updated_at, email FROM users WHERE id = \$1 AND deleted_at IS NULL`).
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := UserFindByIDParams{ID: userID}
	user, err := repo.FindByID(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserGetTimeline_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	userID := uuid.New()
	followedUserID := uuid.New()
	tweetID := uuid.New()
	email := "followed@example.com"
	post := "This is a test tweet"
	now := time.Now()

	// Set up expectations for the complex timeline query
	expectedQuery := `
		WITH followed_users AS \(
			SELECT followed_id 
			FROM follows 
			WHERE follower_id = \$1 AND deleted_at IS NULL
		\),
		tweets_from_followed AS \(
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
				ROW_NUMBER\(\) OVER \(PARTITION BY t.creator_id ORDER BY t.created_at DESC\) AS tweet_rank
			FROM 
				users u
				JOIN tweets t ON u.id = t.creator_id
			WHERE 
				u.id IN \(SELECT followed_id FROM followed_users\)
				AND t.deleted_at IS NULL
				AND u.deleted_at IS NULL
		\)
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
			tweet_rank <= \$2
		ORDER BY 
			tweet_created_at DESC
		LIMIT \$3 OFFSET \$4;
	`

	mock.ExpectQuery(expectedQuery).
		WithArgs(userID, 5, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_id", "email", "user_created_at", "user_updated_at",
			"tweet_id", "tweet_created_at", "tweet_updated_at", "post", "creator_id",
		}).AddRow(
			followedUserID, email, now, now,
			tweetID, now, now, post, followedUserID,
		))

	// Execute the method
	ctx := context.Background()
	params := TimelineParams{
		UserID:            userID,
		TweetsPerFollowed: 5,
		Limit:             10,
		Offset:            0,
	}
	result, err := repo.GetTimeline(ctx, params)

	// Assertions
	require.NoError(t, err)
	require.Len(t, result, 1)

	userWithTweets := result[0]
	assert.Equal(t, followedUserID, userWithTweets.User.ID)
	assert.Equal(t, email, userWithTweets.User.Email)
	require.Len(t, userWithTweets.Tweets, 1)

	tweet := userWithTweets.Tweets[0]
	assert.Equal(t, tweetID, tweet.ID)
	assert.Equal(t, post, tweet.Post)
	assert.Equal(t, followedUserID, tweet.CreatorID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserGetTimeline_EmptyResult(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	userID := uuid.New()

	// Set up expectations for empty result
	expectedQuery := `
		WITH followed_users AS \(
			SELECT followed_id 
			FROM follows 
			WHERE follower_id = \$1 AND deleted_at IS NULL
		\),
		tweets_from_followed AS \(
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
				ROW_NUMBER\(\) OVER \(PARTITION BY t.creator_id ORDER BY t.created_at DESC\) AS tweet_rank
			FROM 
				users u
				JOIN tweets t ON u.id = t.creator_id
			WHERE 
				u.id IN \(SELECT followed_id FROM followed_users\)
				AND t.deleted_at IS NULL
				AND u.deleted_at IS NULL
		\)
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
			tweet_rank <= \$2
		ORDER BY 
			tweet_created_at DESC
		LIMIT \$3 OFFSET \$4;
	`

	mock.ExpectQuery(expectedQuery).
		WithArgs(userID, 5, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_id", "email", "user_created_at", "user_updated_at",
			"tweet_id", "tweet_created_at", "tweet_updated_at", "post", "creator_id",
		})) // No rows added

	ctx := context.Background()
	params := TimelineParams{
		UserID:            userID,
		TweetsPerFollowed: 5,
		Limit:             10,
		Offset:            0,
	}
	result, err := repo.GetTimeline(ctx, params)

	// Assertions
	require.NoError(t, err)
	assert.Empty(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserGetTimeline_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	userID := uuid.New()

	// Set up expectation for database error
	expectedQuery := `
		WITH followed_users AS \(
			SELECT followed_id 
			FROM follows 
			WHERE follower_id = \$1 AND deleted_at IS NULL
		\),
		tweets_from_followed AS \(
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
				ROW_NUMBER\(\) OVER \(PARTITION BY t.creator_id ORDER BY t.created_at DESC\) AS tweet_rank
			FROM 
				users u
				JOIN tweets t ON u.id = t.creator_id
			WHERE 
				u.id IN \(SELECT followed_id FROM followed_users\)
				AND t.deleted_at IS NULL
				AND u.deleted_at IS NULL
		\)
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
			tweet_rank <= \$2
		ORDER BY 
			tweet_created_at DESC
		LIMIT \$3 OFFSET \$4;
	`

	mock.ExpectQuery(expectedQuery).
		WithArgs(userID, 5, 10, 0).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := TimelineParams{
		UserID:            userID,
		TweetsPerFollowed: 5,
		Limit:             10,
		Offset:            0,
	}
	result, err := repo.GetTimeline(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserGetTimeline_MultipleFollowedUsersWithMultipleTweets(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) UserRepository {
		return NewPostgresUserRepository(db)
	})
	defer db.Close()

	// Test data
	userID := uuid.New()
	followedUserID1 := uuid.New()
	followedUserID2 := uuid.New()
	now := time.Now()

	// Generate some unique tweet IDs and times for better testing
	tweetID1_1 := uuid.New()
	tweetID2_1 := uuid.New()

	tweetTime1_1 := now.Add(-5 * time.Minute)
	tweetTime2_1 := now.Add(-3 * time.Minute)

	expectedQuery := `
        WITH followed_users AS \(
            SELECT followed_id 
            FROM follows 
            WHERE follower_id = \$1 AND deleted_at IS NULL
        \),
        tweets_from_followed AS \(
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
                ROW_NUMBER\(\) OVER \(PARTITION BY t.creator_id ORDER BY t.created_at DESC\) AS tweet_rank
            FROM 
                users u
                JOIN tweets t ON u.id = t.creator_id
            WHERE 
                u.id IN \(SELECT followed_id FROM followed_users\)
                AND t.deleted_at IS NULL
                AND u.deleted_at IS NULL
        \)
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
            tweet_rank <= \$2
        ORDER BY 
            tweet_created_at DESC
        LIMIT \$3 OFFSET \$4;
    `

	// Simulate rows for two followed users, each with two tweets
	mock.ExpectQuery(expectedQuery).
		WithArgs(userID, 1, 10, 0). // tweetsPerFollowed = 1
		WillReturnRows(sqlmock.NewRows([]string{
			"user_id", "email", "user_created_at", "user_updated_at",
			"tweet_id", "tweet_created_at", "tweet_updated_at", "post", "creator_id",
		}).
			AddRow(followedUserID1, "followed1@example.com", now, now, tweetID1_1, tweetTime1_1, now, "Tweet 1 by followed 1", followedUserID1).
			AddRow(followedUserID2, "followed2@example.com", now, now, tweetID2_1, tweetTime2_1, now, "Tweet 1 by followed 2", followedUserID2))

	ctx := context.Background()
	params := TimelineParams{
		UserID:            userID,
		TweetsPerFollowed: 1, // Only 1 tweet per followed user
		Limit:             10,
		Offset:            0,
	}
	result, err := repo.GetTimeline(ctx, params)

	require.NoError(t, err)
	require.Len(t, result, 2) // Expect 2 users in the result

	// Assert for followed user 1
	user1 := result[0]
	if user1.User.ID == followedUserID2 {
		user1 = result[1]
	}
	assert.Equal(t, followedUserID1, user1.User.ID)
	assert.Equal(t, "followed1@example.com", user1.User.Email)
	require.Len(t, user1.Tweets, 1)
	assert.Equal(t, tweetID1_1, user1.Tweets[0].ID)

	// Assert for followed user 2
	user2 := result[1]
	if user2.User.ID == followedUserID1 {
		user2 = result[0]
	}
	assert.Equal(t, followedUserID2, user2.User.ID)
	assert.Equal(t, "followed2@example.com", user2.User.Email)
	require.Len(t, user2.Tweets, 1)
	assert.Equal(t, tweetID2_1, user2.Tweets[0].ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
