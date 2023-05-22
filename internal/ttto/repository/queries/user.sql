-- name: CreateUser :one
INSERT INTO users(username, password)
    VALUES (@username, @password)
    RETURNING id;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username=$1;
