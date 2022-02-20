-- name: GetQuipByScore :one
SELECT *
FROM quips
where score_value = $1
  and (not inside_joke or (inside_joke and inside_joke_server_id = $2))
ORDER BY uses, random()
LIMIT 1;

-- name: GetQuipsByCreatedByAccount :many
SELECT *
FROM quips
where created_by_account = $1;

-- name: CreateQuipForScore :one
insert into quips (score_value, quip, inside_joke, inside_joke_server_id, created_by_account)
VALUES ($1, $2, $3, $4, $5)
returning *;

-- name: IncrementQuip :exec
UPDATE quips
SET uses = uses + 1
WHERE id = $1;

-- name: GetQuipsByServerId :many
select *
from quips
where inside_joke_server_id = $1
order by score_value, id;

-- name: DeleteQuipByIdAndServerId :exec
delete
from quips
where id = $1
  and inside_joke_server_id = $2;