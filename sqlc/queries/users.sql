-- name: GetUserById :one
SELECT id, name, email, created_at
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
    RETURNING id, name, email, created_at;
