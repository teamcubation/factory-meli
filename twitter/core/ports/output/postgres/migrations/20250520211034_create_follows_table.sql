-- +goose Up
CREATE TABLE follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    follower_id UUID NOT NULL REFERENCES users(id),
    followed_id UUID NOT NULL REFERENCES users(id)
);

-- Instead of a UNIQUE constraint, we use a partial unique index that only applies to non-deleted rows
CREATE UNIQUE INDEX unique_active_follows ON follows (follower_id, followed_id) WHERE deleted_at IS NULL;

-- Create additional indices to speed up timeline queries
CREATE INDEX idx_follows_follower_id ON follows (follower_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_follows_followed_id ON follows (followed_id) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE follows;