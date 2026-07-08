-- name: CreateFeed :one
INSERT INTO feeds(id, name, url, category, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: FeedExists :one
SELECT EXISTS(
    SELECT 1
    FROM feeds
    WHERE url = $1
);

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetFeedstoFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: IncrementFeedFollowerCount :one
UPDATE feeds
SET feed_followers_count = feed_followers_count + 1
WHERE id = $1
RETURNING *;

-- name: DecrementFeedFollowerCount :one
UPDATE feeds
SET feed_followers_count = feed_followers_count - 1
WHERE id = $1 AND feed_followers_count > 0
RETURNING *;