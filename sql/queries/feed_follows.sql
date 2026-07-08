-- name: FollowFeed :one
INSERT INTO feed_follows (id, user_id, feed_id, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UnfollowFeed :one
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2
RETURNING *;

-- name: GetFollowedFeedsByUser :many
SELECT feeds.* FROM feed_follows
LEFT JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE user_id = $1
ORDER BY feed_follows.created_at DESC;