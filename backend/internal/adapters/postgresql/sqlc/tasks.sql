------------------------ Tasks ------------------------
-- name: GetTaskByID :one
SELECT id,
  owner_id,
  title,
  description,
  due_date,
  completed,
  created_at,
  updated_at
FROM tasks
WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (
    owner_id,
    title,
    description,
    due_date,
    completed
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING id,
  owner_id,
  title,
  description,
  due_date,
  completed,
  created_at,
  updated_at;

-- name: UpdateTask :one
UPDATE tasks
SET title = $2,
  description = $3,
  due_date = $4,
  completed = $5
WHERE id = $1
RETURNING id,
  owner_id,
  title,
  description,
  due_date,
  completed,
  created_at,
  updated_at;

-- name: DeleteTask :execrows
DELETE FROM tasks
WHERE id = $1;