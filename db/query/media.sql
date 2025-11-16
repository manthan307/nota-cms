-- name: CreateMedia :one
INSERT INTO media (key, url, bucket, type, uploaded_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMediaByID :one
SELECT * FROM media
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMedia :many
SELECT * FROM media
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteMedia :exec
UPDATE media
SET deleted_at = now()
WHERE id = $1;

-- name: UpdateMedia :one
UPDATE media
SET key = $2, url = $3, bucket = $4, type = $5, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: GetMediaByURL :one
SELECT * FROM media
WHERE url = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;