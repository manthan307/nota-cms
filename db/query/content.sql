-- name: CreateContent :one
INSERT INTO contents (schema_id, data, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetContentsBySchema :many
SELECT * FROM contents
WHERE schema_id = $1
AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteContent :exec
UPDATE contents
SET deleted_at = now()
WHERE id = $1;
