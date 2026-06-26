-- name: GetChirpsByAuthor :many

SELECT * FROM chirps
where user_id = $1
ORDER BY created_at ASC;