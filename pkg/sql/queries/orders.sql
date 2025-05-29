-- name: CreateOrder :one
INSERT INTO orders (id, userid, price, quantity, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
