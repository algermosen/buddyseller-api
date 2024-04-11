-- name: CreateUser :one
INSERT INTO users(name, code, email, password)
VALUES ($1, $2, $3, $4) RETURNING id;

-- name: ListUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUser :execrows
UPDATE users
SET 
    name = $2,
    code = $3,
    email = $4
WHERE id = $1;

-- name: UpdatePassword :execrows
UPDATE users
SET 
    password = $2
WHERE id = $1;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1;

-- name: GetUserCredentials :one
SELECT id, email, password FROM users
WHERE code = $1;