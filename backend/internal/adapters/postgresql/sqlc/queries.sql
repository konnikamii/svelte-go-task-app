-- name: ListProducts :many
SELECT * FROM tasks;

-- name: GetProductById :one
SELECT * FROM tasks WHERE id = $1;