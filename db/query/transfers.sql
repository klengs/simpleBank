-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetTransferByFromToIDS :many
SELECT *
FROM transfers
WHERE
    from_account_id = $1
    AND to_account_id = $2
LIMIT $3
OFFSET $4;