-- name: CreateUser :execresult
INSERT INTO users (username, password_hash) VALUES (?, ?);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? LIMIT 1;
