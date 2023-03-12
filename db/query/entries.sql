-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetEntryByID :one
SELECT *
FROM entries
WHERE id = $1;

-- name: GetEntriesByAccountID :many
SELECT *
FROM entries
WHERE account_id = $1
LIMIT $2
OFFSET $3;

-- name: GetAllEntries :many
SELECT *
FROM entries
LIMIT $1
OFFSET $2;