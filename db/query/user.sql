-- name: CreateUser :one
INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY id;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1 FROM users
    WHERE id = $1
    AND deleted_at IS NULL
) AS exists;

-- name: AdminExists :one
SELECT EXISTS (
    SELECT 1 FROM users
    WHERE role = 'admin'
    AND deleted_at IS NULL
) AS exists;