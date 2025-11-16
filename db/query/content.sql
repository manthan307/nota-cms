-- name: CreateContent :one
INSERT INTO contents (schema_id, data, created_by,published)
VALUES ($1, $2, $3,$4)
RETURNING *;

-- name: GetContentsBySchema :many
SELECT * FROM contents
WHERE schema_id = $1
AND deleted_at IS NULL
AND published = $2
ORDER BY created_at DESC;

-- name: GetAllContentsBySchema :many
SELECT * FROM contents
WHERE schema_id = $1
AND deleted_at IS NULL
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
SET 
  data = $2,
  published = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: GetAllContents :many
SELECT * FROM contents
WHERE deleted_at IS NULL
ORDER BY created_at DESC;
