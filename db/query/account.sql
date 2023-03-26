-- name: CreateAccount :one
INSERT INTO account (
    owner,
    balance,
    currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccountByID :one
SELECT *
FROM account
WHERE id = $1 
LIMIT 1;

-- name: GetAllAccounts :many
SELECT *
FROM account
WHERE owner = $1
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE
    account
SET 
    balance = $1,
    currency = $2,
    owner = $3
WHERE
    id = $4
RETURNING *;

-- name: UpdateBalance :one
UPDATE
    account
SET 
    balance = balance + $1
WHERE
    id = $2
RETURNING *;

-- name: DeleteAccount :one
DELETE FROM account
WHERE id = $1
RETURNING *;