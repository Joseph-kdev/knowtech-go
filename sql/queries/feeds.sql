-- name: CreateFeed :one
INSERT INTO feeds(id, name, url, category, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: FeedExists :one
SELECT EXISTS(
    SELECT 1
    FROM feeds
    WHERE url = ?
);

-- name: GetAllFeeds :many
SELECT * FROM feeds;