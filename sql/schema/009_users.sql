-- +goose Up
ALTER TABLE users RENAME COLUMN name TO email;

-- +goose Down
ALTER TABLE users RENAME COLUMN email TO name;