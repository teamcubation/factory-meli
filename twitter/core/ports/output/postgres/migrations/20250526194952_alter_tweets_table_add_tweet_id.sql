-- +goose Up
-- Add parent_id column for tweet threading
ALTER TABLE tweets ADD COLUMN parent_id UUID REFERENCES tweets(id);

-- Index for efficient thread traversal: find all replies to a specific tweet
CREATE INDEX idx_tweets_parent_id ON tweets (parent_id) WHERE deleted_at IS NULL AND parent_id IS NOT NULL;

-- Composite index for thread queries with ordering by creation time
CREATE INDEX idx_tweets_parent_created_active ON tweets (parent_id, created_at ASC) WHERE deleted_at IS NULL AND parent_id IS NOT NULL;

-- Index to efficiently find root tweets (original tweets without parent)
CREATE INDEX idx_tweets_root_tweets ON tweets (creator_id, created_at DESC) WHERE deleted_at IS NULL AND parent_id IS NULL;

-- +goose Down
-- Drop indexes first (in this case, we need to do that because the table will still exist if this migration is reverted)
DROP INDEX IF EXISTS idx_tweets_root_tweets;
DROP INDEX IF EXISTS idx_tweets_parent_created_active;
DROP INDEX IF EXISTS idx_tweets_parent_id;

-- Drop the column
ALTER TABLE tweets DROP COLUMN parent_id;