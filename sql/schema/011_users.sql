-- +goose Up
ALTER TABLE feed_follows DROP CONSTRAINT IF EXISTS feed_follows_user_id_fkey;
ALTER TABLE users ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE feed_follows ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;
ALTER TABLE feed_follows
ADD CONSTRAINT feed_follows_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE feed_follows DROP CONSTRAINT IF EXISTS feed_follows_user_id_fkey;
ALTER TABLE feed_follows ALTER COLUMN user_id TYPE UUID USING user_id::UUID;
ALTER TABLE users ALTER COLUMN id TYPE UUID USING id::UUID;
ALTER TABLE feed_follows
ADD CONSTRAINT feed_follows_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;