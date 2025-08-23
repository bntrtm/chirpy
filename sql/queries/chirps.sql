-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteChirpByID :exec
DELETE
FROM chirps
WHERE chirps.id = $1;

-- name: GetRecentChirps :many
SELECT *
FROM chirps
WHERE ($1 = '') OR (user_id = $1::uuid)
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT *
FROM chirps
WHERE chirps.id = $1;