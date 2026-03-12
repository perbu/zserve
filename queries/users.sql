-- name: GetUser :one
SELECT * FROM zendesk.users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM zendesk.users
ORDER BY name ASC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: CountUsers :one
SELECT count(*) FROM zendesk.users;

-- name: GetUserNames :many
SELECT id, name FROM zendesk.users WHERE id = ANY(@ids::bigint[]);
