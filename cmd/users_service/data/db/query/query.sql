-- name: CreateUser :one
INSERT INTO users(username, password_hash)
VALUES(?, ?)
RETURNING id, username; 

-- name: GetUserByUsername :one
SELECT *
FROM users
where username = ?;