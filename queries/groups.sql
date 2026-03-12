-- name: GetGroup :one
SELECT * FROM zendesk.groups WHERE id = $1;

-- name: ListGroups :many
SELECT * FROM zendesk.groups
ORDER BY name ASC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: CountGroups :one
SELECT count(*) FROM zendesk.groups;

-- name: GetGroupNames :many
SELECT id, name FROM zendesk.groups WHERE id = ANY(@ids::bigint[]);
