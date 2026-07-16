-- name: AddToBookmarks :one
INSERT INTO bookmarks (id, user_id, feed_id, title, url, description, published_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (user_id, feed_id, url) DO UPDATE SET updated_at = EXCLUDED.updated_at
RETURNING id, user_id, feed_id, title, url, description, published_at, created_at, updated_at;

-- name: GetBookmarksByUser :many
SELECT 
    b.id,
    b.title,
    b.url,
    b.description,
    b.published_at,
    b.created_at,
    b.updated_at,
    f.id AS feed_id,
    f.name AS feed_name,
    f.url AS feed_url
FROM bookmarks b
LEFT JOIN feeds f ON b.feed_id = f.id
WHERE b.user_id = $1
ORDER BY b.published_at DESC;

-- name: DeleteBookmark :one
DELETE FROM bookmarks
WHERE id = $1 AND user_id = $2
RETURNING *;