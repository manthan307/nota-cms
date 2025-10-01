-- name: CreateContent :one
INSERT INTO contents (schema_id, data, created_by,published)
VALUES ($1, $2, $3,$4)
RETURNING *;

-- name: GetContentsBySchema :many
SELECT * FROM contents
WHERE schema_id = $1
AND deleted_at IS NULL
AND published = true
ORDER BY created_at DESC;

-- name: DeleteContent :exec
UPDATE contents
SET deleted_at = now()
WHERE id = $1;

-- name: GetContentByID :one
SELECT * FROM contents
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateContent :one
UPDATE contents
SET data = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: GetAllContents :many
SELECT * FROM contents
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: SetContentPublished :one
UPDATE contents
SET published = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
