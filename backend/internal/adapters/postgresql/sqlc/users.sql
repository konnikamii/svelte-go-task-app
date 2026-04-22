------------------------ Users ------------------------
-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, PASSWORD)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2,
  email = $3,
  PASSWORD = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1;