-- name: CreateTransaction :one
INSERT INTO transactions(
    account_id, amount
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE account_id = $1
ORDER by id
LIMIT $2
OFFSET $3;