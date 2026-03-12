-- name: GetOrganization :one
SELECT * FROM zendesk.organizations WHERE id = $1;

-- name: ListOrganizations :many
SELECT * FROM zendesk.organizations
ORDER BY name ASC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: CountOrganizations :one
SELECT count(*) FROM zendesk.organizations;

-- name: GetOrganizationNames :many
SELECT id, name FROM zendesk.organizations WHERE id = ANY(@ids::bigint[]);
