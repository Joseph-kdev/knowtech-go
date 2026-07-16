-- +goose Up
ALTER TABLE bookmarks ADD CONSTRAINT bookmarks_unique UNIQUE (user_id, feed_id, url);

-- +goose Down
ALTER TABLE bookmarks DROP CONSTRAINT bookmarks_unique;