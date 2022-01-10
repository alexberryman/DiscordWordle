-- name: GetResponseByScore :one
SELECT *
FROM responses
where score_value = $1
  and (not inside_joke or (inside_joke and inside_joke_server_id = $2))
ORDER BY random()
LIMIT 1;

-- name: GetResponsesByCreatedByAccount :many
SELECT *
FROM responses
where created_by_account = $1;

-- name: CreateResponseForScore :one
insert into responses (score_value, response, inside_joke, inside_joke_server_id, created_by_account)
VALUES ($1, $2, $3, $4, $5) returning *;