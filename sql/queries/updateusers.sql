-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, email, hashed_password, created_at, updated_at, is_chirpy_red;