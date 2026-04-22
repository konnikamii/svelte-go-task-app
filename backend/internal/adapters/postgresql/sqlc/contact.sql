------------------------ Contact ------------------------
-- name: CreateContactRequest :one
INSERT INTO contact_requests (email, title, message)
VALUES ($1, $2, $3)
RETURNING id,
  email,
  title,
  message,
  created_at,
  updated_at;