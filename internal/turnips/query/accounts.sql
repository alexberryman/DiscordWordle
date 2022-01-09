-- name: GetAccount :one
SELECT *
FROM accounts
WHERE discord_id = $1
LIMIT 1;

-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY discord_id;

-- name: CountAccountsByDiscordId :one
SELECT count(*)
FROM accounts
where discord_id = $1;

-- name: CreateAccount :one
INSERT INTO accounts (discord_id)
VALUES ($1)
RETURNING *;

-- name: UpdateTimeZone :one
UPDATE accounts
set time_zone = $2
where discord_id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE discord_id = $1;