-- name: AddPostsToDatabase :many
INSERT INTO posts (id, feed_id, title, url, description, published_at, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsByFeed :many
SELECT 
    f.id AS feed_id,
    f.name AS feed_name,
    f.url AS feed_url,
    json_agg(
        json_build_object(
            'id', p.id,
            'title', p.title,
            'url', p.url,
            'description', p.description,
            'published_at', p.published_at,
            'created_at', p.created_at,
            'updated_at', p.updated_at
        ) ORDER BY p.published_at DESC
    ) AS posts
FROM feeds f
JOIN posts p ON p.feed_id = f.id
GROUP BY f.id, f.name, f.url
ORDER BY f.name;

-- -- name: DeleteStalePosts :exec
-- DELETE FROM posts
-- WHERE published_at < NOW() - INTERVAL "12 hours";

