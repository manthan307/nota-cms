-- name: CreateSchema :one
INSERT INTO schemas (name, definition, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSchemaByID :one
SELECT * FROM schemas
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetSchemaByName :one
SELECT * FROM schemas
WHERE name = $1 AND deleted_at IS NULL;

-- name: ListSchemas :many
SELECT * FROM schemas
WHERE deleted_at IS NULL
ORDER BY id;

-- name: DeleteSchema :exec
DELETE FROM schemas
WHERE id = $1;

-- name: UpdateSchema :one
UPDATE schemas
SET name = $2, definition = $3, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;