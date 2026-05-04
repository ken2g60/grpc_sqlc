-- name: CreateUser :execresult
INSERT INTO users (first_name, last_name, phone_number)
VALUES ($1, $2, $3) RETURNING id, first_name, last_name, phone_number;

-- name: UpdateUser :exec
UPDATE users SET first_name = $1, last_name = $2, phone_number = $3
WHERE id = $4 RETURNING id, first_name, last_name, phone_number;

-- name: GetUser :one
SELECT id, first_name, last_name, phone_number FROM users WHERE id = $1;

-- name: DeactivateUser :exec
DELETE FROM users WHERE id = $1 RETURNING id, first_name, last_name, phone_number;