-- name: GetNickname :one
SELECT *
FROM nicknames
WHERE discord_id = $1
and server_id = $2
LIMIT 1;

-- name: ListNicknames :many
SELECT *
FROM nicknames
ORDER BY discord_id;

-- name: CountNicknameByDiscordId :one
SELECT count(*)
FROM nicknames
where discord_id = $1
  and server_id = $2;

-- name: CreateNickname :one
INSERT INTO nicknames (discord_id, server_id, nickname)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateNickname :one
UPDATE nicknames
set nickname = $2
where discord_id = $1
  and server_id = $3
RETURNING *;

-- name: DeleteNickname :exec
DELETE
FROM nicknames
WHERE discord_id = $1;