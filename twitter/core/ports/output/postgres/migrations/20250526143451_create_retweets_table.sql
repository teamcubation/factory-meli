-- +goose Up
CREATE TABLE retweets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    user_id UUID NOT NULL REFERENCES users(id),
    tweet_id UUID NOT NULL REFERENCES tweets(id)
);

-- Instead of a UNIQUE constraint, we use a partial unique index that only applies to non-deleted rows
CREATE UNIQUE INDEX idx_user_tweet_retweet_unique ON retweets (user_id, tweet_id) WHERE deleted_at IS NULL;

-- Index for finding all retweets by a user
CREATE INDEX idx_retweets_user_id ON retweets (user_id) WHERE deleted_at IS NULL;

-- Index for finding all retweets of a specific tweet (for counting retweets)
CREATE INDEX idx_retweets_tweet_id ON retweets (tweet_id) WHERE deleted_at IS NULL;

-- Composite index for user timeline queries (user's retweets ordered by creation date)
CREATE INDEX idx_retweets_user_created_active ON retweets (user_id, created_at DESC) WHERE deleted_at IS NULL;

-- Composite index for tweet engagement queries (tweet's retweets ordered by date)
CREATE INDEX idx_retweets_tweet_created_active ON retweets (tweet_id, created_at DESC) WHERE deleted_at IS NULL;

-- Index for soft delete filtering
CREATE INDEX idx_retweets_deleted_at ON retweets (deleted_at) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE retweets;