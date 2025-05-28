-- +goose Up
ALTER TABLE users 
ADD COLUMN name VARCHAR(100),
ADD COLUMN bio VARCHAR(280),
ADD COLUMN avatar_url TEXT;

-- +goose Down
ALTER TABLE users 
DROP COLUMN name,
DROP COLUMN bio,
DROP COLUMN avatar_url;