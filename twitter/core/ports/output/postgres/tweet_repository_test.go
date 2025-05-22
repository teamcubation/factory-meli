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

func TestTweetSave_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	// Test data
	post := "My first tweet!"
	creatorID := uuid.New()
	tweetID := uuid.New()
	now := time.Now()

	// Set up expectations
	mock.ExpectQuery(`INSERT INTO tweets \(post, creator_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, post, creator_id;`).
		WithArgs(post, creatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "post", "creator_id"}).
			AddRow(tweetID, now, now, post, creatorID))

	// Execute the method
	ctx := context.Background()
	params := TweetSaveParams{Post: post, CreatorID: creatorID}
	tweet, err := repo.Save(ctx, params)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, tweetID, tweet.ID)
	assert.Equal(t, post, tweet.Post)
	assert.Equal(t, creatorID, tweet.CreatorID)
	assert.WithinDuration(t, now, tweet.CreatedAt, time.Second)
	assert.WithinDuration(t, now, tweet.UpdatedAt, time.Second)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetSave_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	post := "Error tweet"
	creatorID := uuid.New()

	// Set up expectation for database error
	mock.ExpectQuery(`INSERT INTO tweets \(post, creator_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, post, creator_id;`).
		WithArgs(post, creatorID).
		WillReturnError(sql.ErrConnDone) // Simulate connection error

	ctx := context.Background()
	params := TweetSaveParams{Post: post, CreatorID: creatorID}
	tweet, err := repo.Save(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, tweet)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetSave_ForeignKeyViolation(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	post := "Tweet by non-existent user"
	creatorID := uuid.New()

	expectedErrorMsg := `pq: insert or update on table "tweets" violates foreign key constraint "tweets_creator_id_fkey"`
	mock.ExpectQuery(`INSERT INTO tweets \(post, creator_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, post, creator_id;`).
		WithArgs(post, creatorID).
		WillReturnError(errors.New(expectedErrorMsg))

	ctx := context.Background()
	params := TweetSaveParams{Post: post, CreatorID: creatorID}
	tweet, err := repo.Save(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, tweet)
	assert.Contains(t, err.Error(), "violates foreign key constraint")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetFindByID_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	// Test data
	tweetID := uuid.New()
	creatorID := uuid.New()
	post := "Found tweet"
	now := time.Now()

	// Set up expectations
	mock.ExpectQuery(`SELECT id, created_at, updated_at, post, creator_id FROM tweets WHERE id = \$1 AND deleted_at IS NULL;`).
		WithArgs(tweetID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "post", "creator_id"}).
			AddRow(tweetID, now, now, post, creatorID))

	// Execute the method
	ctx := context.Background()
	params := TweetFindByIDParams{ID: tweetID}
	tweet, err := repo.FindByID(ctx, params)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, tweetID, tweet.ID)
	assert.Equal(t, post, tweet.Post)
	assert.Equal(t, creatorID, tweet.CreatorID)
	assert.WithinDuration(t, now, tweet.CreatedAt, time.Second)
	assert.WithinDuration(t, now, tweet.UpdatedAt, time.Second)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetFindByID_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	tweetID := uuid.New()

	// Set up expectation for no rows found
	mock.ExpectQuery(`SELECT id, created_at, updated_at, post, creator_id FROM tweets WHERE id = \$1 AND deleted_at IS NULL;`).
		WithArgs(tweetID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	params := TweetFindByIDParams{ID: tweetID}
	tweet, err := repo.FindByID(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, tweet)
	assert.Equal(t, sql.ErrNoRows, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetFindByID_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	tweetID := uuid.New()

	// Set up expectation for database connection error
	mock.ExpectQuery(`SELECT id, created_at, updated_at, post, creator_id FROM tweets WHERE id = \$1 AND deleted_at IS NULL;`).
		WithArgs(tweetID).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := TweetFindByIDParams{ID: tweetID}
	tweet, err := repo.FindByID(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, tweet)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetFindAllTweetsByUserId_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	creatorID := uuid.New()
	now := time.Now()

	tweetID1 := uuid.New()
	tweetPost1 := "First tweet by user"
	tweetCreatedAt1 := now.Add(-5 * time.Minute)

	tweetID2 := uuid.New()
	tweetPost2 := "Second tweet by user"
	tweetCreatedAt2 := now.Add(-10 * time.Minute)

	expectedQuery := `
        SELECT id, created_at, updated_at, post, creator_id
            FROM tweets WHERE creator_id = \$1
            ORDER BY created_at DESC LIMIT \$2 OFFSET \$3;
    `

	// Simulate two tweets for the user
	mock.ExpectQuery(expectedQuery).
		WithArgs(creatorID, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "post", "creator_id"}).
			AddRow(tweetID1, tweetCreatedAt1, now, tweetPost1, creatorID). // More recent
			AddRow(tweetID2, tweetCreatedAt2, now, tweetPost2, creatorID)) // Less recent

	ctx := context.Background()
	params := TweetFindAllTweetsByUserId{
		CreatorID: creatorID,
		Limit:     10,
		Offset:    0,
	}
	tweets, err := repo.FindAllTweetsByUserId(ctx, params)

	require.NoError(t, err)
	require.Len(t, tweets, 2)

	// Assert the order (most recent first due to ORDER BY created_at DESC)
	assert.Equal(t, tweetID1, tweets[0].ID)
	assert.Equal(t, tweetPost1, tweets[0].Post)
	assert.Equal(t, creatorID, tweets[0].CreatorID)
	assert.WithinDuration(t, tweetCreatedAt1, tweets[0].CreatedAt, time.Second)

	assert.Equal(t, tweetID2, tweets[1].ID)
	assert.Equal(t, tweetPost2, tweets[1].Post)
	assert.Equal(t, creatorID, tweets[1].CreatorID)
	assert.WithinDuration(t, tweetCreatedAt2, tweets[1].CreatedAt, time.Second)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetFindAllTweetsByUserId_NoTweetsFound(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	creatorID := uuid.New()

	expectedQuery := `
        SELECT id, created_at, updated_at, post, creator_id
            FROM tweets WHERE creator_id = \$1
            ORDER BY created_at DESC LIMIT \$2 OFFSET \$3;
    `

	// Return no rows
	mock.ExpectQuery(expectedQuery).
		WithArgs(creatorID, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "post", "creator_id"}))

	ctx := context.Background()
	params := TweetFindAllTweetsByUserId{
		CreatorID: creatorID,
		Limit:     10,
		Offset:    0,
	}
	tweets, err := repo.FindAllTweetsByUserId(ctx, params)

	require.NoError(t, err)
	assert.Empty(t, tweets) // Expect an empty slice

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetFindAllTweetsByUserId_WithLimitAndOffset(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	creatorID := uuid.New()
	now := time.Now()

	// Simulate 4 tweets, ordered by creation time DESC
	tweetID1 := uuid.New()
	tweetTime1 := now.Add(-2 * time.Minute)
	tweetID2 := uuid.New()
	tweetTime2 := now.Add(-3 * time.Minute)

	expectedQuery := `
        SELECT id, created_at, updated_at, post, creator_id
            FROM tweets WHERE creator_id = \$1
            ORDER BY created_at DESC LIMIT \$2 OFFSET \$3;
    `

	// The order should reflect `ORDER BY created_at DESC`.
	mock.ExpectQuery(expectedQuery).
		WithArgs(creatorID, 2, 1). // Limit 2, Offset 1
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "post", "creator_id"}).
			AddRow(tweetID1, tweetTime1, now, "Tweet 2", creatorID). // This is the 1st row (after offset)
			AddRow(tweetID2, tweetTime2, now, "Tweet 3", creatorID)) // This is the 2nd row

	ctx := context.Background()
	params := TweetFindAllTweetsByUserId{
		CreatorID: creatorID,
		Limit:     2,
		Offset:    1,
	}
	tweets, err := repo.FindAllTweetsByUserId(ctx, params)

	require.NoError(t, err)
	require.Len(t, tweets, 2)

	assert.Equal(t, tweetID1, tweets[0].ID)
	assert.Equal(t, tweetID2, tweets[1].ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetFindAllTweetsByUserId_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	creatorID := uuid.New()

	expectedQuery := `
        SELECT id, created_at, updated_at, post, creator_id
            FROM tweets WHERE creator_id = \$1
            ORDER BY created_at DESC LIMIT \$2 OFFSET \$3;
    `

	// Set up expectation for database error
	mock.ExpectQuery(expectedQuery).
		WithArgs(creatorID, 10, 0).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := TweetFindAllTweetsByUserId{
		CreatorID: creatorID,
		Limit:     10,
		Offset:    0,
	}
	tweets, err := repo.FindAllTweetsByUserId(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, tweets)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTweetFindAllTweetsByUserId_ScanningError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) TweetRepository {
		return NewPostgresTweetRepository(db)
	})
	defer db.Close()

	creatorID := uuid.New()
	now := time.Now()

	expectedQuery := `
        SELECT id, created_at, updated_at, post, creator_id
            FROM tweets WHERE creator_id = \$1
            ORDER BY created_at DESC LIMIT \$2 OFFSET \$3;
    `

	// Provide a row with a type mismatch
	mock.ExpectQuery(expectedQuery).
		WithArgs(creatorID, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "post", "creator_id"}).
			AddRow("not-a-uuid", now, now, "Invalid tweet", creatorID))

	ctx := context.Background()
	params := TweetFindAllTweetsByUserId{
		CreatorID: creatorID,
		Limit:     10,
		Offset:    0,
	}
	tweets, err := repo.FindAllTweetsByUserId(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, tweets)
	assert.Contains(t, err.Error(), `Scan error on column index 0, name "id": Scan: invalid UUID length: 10`)
	assert.NoError(t, mock.ExpectationsWereMet())
}
