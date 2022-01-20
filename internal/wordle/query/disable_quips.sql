-- name: CheckIfServerHasDisabledQuips :one
SELECT server_id
from disable_quips
where server_id = $1;

-- name: DisableQuipsForServer :exec
insert into disable_quips (server_id) values ($1);

-- name: EnableQuipsForServer :exec
delete from disable_quips where server_id = $1;