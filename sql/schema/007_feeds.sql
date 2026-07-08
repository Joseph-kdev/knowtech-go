-- +goose Up
ALTER TABLE feeds ADD COLUMN feed_followers_count INTEGER DEFAULT 0;

-- +goose Down
ALTER TABLE feeds DROP COLUMN feed_followers_count;