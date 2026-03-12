-- name: GetTicket :one
SELECT * FROM zendesk.tickets WHERE id = $1;

-- name: CountTickets :one
SELECT count(*) FROM zendesk.tickets t
WHERE
    (sqlc.narg('status')::varchar IS NULL OR t.status = sqlc.narg('status')) AND
    (sqlc.narg('ticket_type')::varchar IS NULL OR t.type = sqlc.narg('ticket_type')) AND
    (sqlc.narg('priority')::varchar IS NULL OR t.priority = sqlc.narg('priority')) AND
    (sqlc.narg('requester_id')::bigint IS NULL OR t.requester_id = sqlc.narg('requester_id')) AND
    (sqlc.narg('assignee_id')::bigint IS NULL OR t.assignee_id = sqlc.narg('assignee_id')) AND
    (sqlc.narg('submitter_id')::bigint IS NULL OR t.submitter_id = sqlc.narg('submitter_id')) AND
    (sqlc.narg('organization_id')::bigint IS NULL OR t.organization_id = sqlc.narg('organization_id')) AND
    (sqlc.narg('group_id')::bigint IS NULL OR t.group_id = sqlc.narg('group_id')) AND
    (sqlc.narg('brand_id')::bigint IS NULL OR t.brand_id = sqlc.narg('brand_id')) AND
    (sqlc.narg('search')::varchar IS NULL OR t.subject ILIKE '%' || sqlc.narg('search') || '%') AND
    (sqlc.narg('created_after')::timestamptz IS NULL OR t.created_at >= sqlc.narg('created_after')) AND
    (sqlc.narg('created_before')::timestamptz IS NULL OR t.created_at <= sqlc.narg('created_before')) AND
    (sqlc.narg('tag')::varchar IS NULL OR EXISTS (
        SELECT 1 FROM zendesk.tickets__tags tag
        WHERE tag._dlt_root_id = t._dlt_id AND tag.value = sqlc.narg('tag')
    ));

-- name: ListTicketsByCreatedAtDesc :many
SELECT * FROM zendesk.tickets t
WHERE
    (sqlc.narg('status')::varchar IS NULL OR t.status = sqlc.narg('status')) AND
    (sqlc.narg('ticket_type')::varchar IS NULL OR t.type = sqlc.narg('ticket_type')) AND
    (sqlc.narg('priority')::varchar IS NULL OR t.priority = sqlc.narg('priority')) AND
    (sqlc.narg('requester_id')::bigint IS NULL OR t.requester_id = sqlc.narg('requester_id')) AND
    (sqlc.narg('assignee_id')::bigint IS NULL OR t.assignee_id = sqlc.narg('assignee_id')) AND
    (sqlc.narg('submitter_id')::bigint IS NULL OR t.submitter_id = sqlc.narg('submitter_id')) AND
    (sqlc.narg('organization_id')::bigint IS NULL OR t.organization_id = sqlc.narg('organization_id')) AND
    (sqlc.narg('group_id')::bigint IS NULL OR t.group_id = sqlc.narg('group_id')) AND
    (sqlc.narg('brand_id')::bigint IS NULL OR t.brand_id = sqlc.narg('brand_id')) AND
    (sqlc.narg('search')::varchar IS NULL OR t.subject ILIKE '%' || sqlc.narg('search') || '%') AND
    (sqlc.narg('created_after')::timestamptz IS NULL OR t.created_at >= sqlc.narg('created_after')) AND
    (sqlc.narg('created_before')::timestamptz IS NULL OR t.created_at <= sqlc.narg('created_before')) AND
    (sqlc.narg('tag')::varchar IS NULL OR EXISTS (
        SELECT 1 FROM zendesk.tickets__tags tag
        WHERE tag._dlt_root_id = t._dlt_id AND tag.value = sqlc.narg('tag')
    ))
ORDER BY t.created_at DESC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: ListTicketsByCreatedAtAsc :many
SELECT * FROM zendesk.tickets t
WHERE
    (sqlc.narg('status')::varchar IS NULL OR t.status = sqlc.narg('status')) AND
    (sqlc.narg('ticket_type')::varchar IS NULL OR t.type = sqlc.narg('ticket_type')) AND
    (sqlc.narg('priority')::varchar IS NULL OR t.priority = sqlc.narg('priority')) AND
    (sqlc.narg('requester_id')::bigint IS NULL OR t.requester_id = sqlc.narg('requester_id')) AND
    (sqlc.narg('assignee_id')::bigint IS NULL OR t.assignee_id = sqlc.narg('assignee_id')) AND
    (sqlc.narg('submitter_id')::bigint IS NULL OR t.submitter_id = sqlc.narg('submitter_id')) AND
    (sqlc.narg('organization_id')::bigint IS NULL OR t.organization_id = sqlc.narg('organization_id')) AND
    (sqlc.narg('group_id')::bigint IS NULL OR t.group_id = sqlc.narg('group_id')) AND
    (sqlc.narg('brand_id')::bigint IS NULL OR t.brand_id = sqlc.narg('brand_id')) AND
    (sqlc.narg('search')::varchar IS NULL OR t.subject ILIKE '%' || sqlc.narg('search') || '%') AND
    (sqlc.narg('created_after')::timestamptz IS NULL OR t.created_at >= sqlc.narg('created_after')) AND
    (sqlc.narg('created_before')::timestamptz IS NULL OR t.created_at <= sqlc.narg('created_before')) AND
    (sqlc.narg('tag')::varchar IS NULL OR EXISTS (
        SELECT 1 FROM zendesk.tickets__tags tag
        WHERE tag._dlt_root_id = t._dlt_id AND tag.value = sqlc.narg('tag')
    ))
ORDER BY t.created_at ASC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: ListTicketsByUpdatedAtDesc :many
SELECT * FROM zendesk.tickets t
WHERE
    (sqlc.narg('status')::varchar IS NULL OR t.status = sqlc.narg('status')) AND
    (sqlc.narg('ticket_type')::varchar IS NULL OR t.type = sqlc.narg('ticket_type')) AND
    (sqlc.narg('priority')::varchar IS NULL OR t.priority = sqlc.narg('priority')) AND
    (sqlc.narg('requester_id')::bigint IS NULL OR t.requester_id = sqlc.narg('requester_id')) AND
    (sqlc.narg('assignee_id')::bigint IS NULL OR t.assignee_id = sqlc.narg('assignee_id')) AND
    (sqlc.narg('submitter_id')::bigint IS NULL OR t.submitter_id = sqlc.narg('submitter_id')) AND
    (sqlc.narg('organization_id')::bigint IS NULL OR t.organization_id = sqlc.narg('organization_id')) AND
    (sqlc.narg('group_id')::bigint IS NULL OR t.group_id = sqlc.narg('group_id')) AND
    (sqlc.narg('brand_id')::bigint IS NULL OR t.brand_id = sqlc.narg('brand_id')) AND
    (sqlc.narg('search')::varchar IS NULL OR t.subject ILIKE '%' || sqlc.narg('search') || '%') AND
    (sqlc.narg('created_after')::timestamptz IS NULL OR t.created_at >= sqlc.narg('created_after')) AND
    (sqlc.narg('created_before')::timestamptz IS NULL OR t.created_at <= sqlc.narg('created_before')) AND
    (sqlc.narg('tag')::varchar IS NULL OR EXISTS (
        SELECT 1 FROM zendesk.tickets__tags tag
        WHERE tag._dlt_root_id = t._dlt_id AND tag.value = sqlc.narg('tag')
    ))
ORDER BY t.updated_at DESC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: ListTicketsByUpdatedAtAsc :many
SELECT * FROM zendesk.tickets t
WHERE
    (sqlc.narg('status')::varchar IS NULL OR t.status = sqlc.narg('status')) AND
    (sqlc.narg('ticket_type')::varchar IS NULL OR t.type = sqlc.narg('ticket_type')) AND
    (sqlc.narg('priority')::varchar IS NULL OR t.priority = sqlc.narg('priority')) AND
    (sqlc.narg('requester_id')::bigint IS NULL OR t.requester_id = sqlc.narg('requester_id')) AND
    (sqlc.narg('assignee_id')::bigint IS NULL OR t.assignee_id = sqlc.narg('assignee_id')) AND
    (sqlc.narg('submitter_id')::bigint IS NULL OR t.submitter_id = sqlc.narg('submitter_id')) AND
    (sqlc.narg('organization_id')::bigint IS NULL OR t.organization_id = sqlc.narg('organization_id')) AND
    (sqlc.narg('group_id')::bigint IS NULL OR t.group_id = sqlc.narg('group_id')) AND
    (sqlc.narg('brand_id')::bigint IS NULL OR t.brand_id = sqlc.narg('brand_id')) AND
    (sqlc.narg('search')::varchar IS NULL OR t.subject ILIKE '%' || sqlc.narg('search') || '%') AND
    (sqlc.narg('created_after')::timestamptz IS NULL OR t.created_at >= sqlc.narg('created_after')) AND
    (sqlc.narg('created_before')::timestamptz IS NULL OR t.created_at <= sqlc.narg('created_before')) AND
    (sqlc.narg('tag')::varchar IS NULL OR EXISTS (
        SELECT 1 FROM zendesk.tickets__tags tag
        WHERE tag._dlt_root_id = t._dlt_id AND tag.value = sqlc.narg('tag')
    ))
ORDER BY t.updated_at ASC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: GetTicketTags :many
SELECT tag.value FROM zendesk.tickets__tags tag
JOIN zendesk.tickets t ON tag._dlt_root_id = t._dlt_id
WHERE t.id = $1
ORDER BY tag._dlt_list_idx;

-- name: GetTicketCollaborators :many
SELECT c.value FROM zendesk.tickets__collaborator_ids c
JOIN zendesk.tickets t ON c._dlt_root_id = t._dlt_id
WHERE t.id = $1
ORDER BY c._dlt_list_idx;

-- name: GetTicketFollowers :many
SELECT f.value FROM zendesk.tickets__follower_ids f
JOIN zendesk.tickets t ON f._dlt_root_id = t._dlt_id
WHERE t.id = $1
ORDER BY f._dlt_list_idx;

-- name: GetTicketEmailCCs :many
SELECT e.value FROM zendesk.tickets__email_cc_ids e
JOIN zendesk.tickets t ON e._dlt_root_id = t._dlt_id
WHERE t.id = $1
ORDER BY e._dlt_list_idx;

-- name: GetTicketFollowups :many
SELECT fu.value FROM zendesk.tickets__followup_ids fu
JOIN zendesk.tickets t ON fu._dlt_root_id = t._dlt_id
WHERE t.id = $1
ORDER BY fu._dlt_list_idx;

-- name: GetTicketComments :many
SELECT id, type, author_id, plain_body, public, created_at, via__channel
FROM zendesk.ticket_comments
WHERE ticket_id = $1
ORDER BY created_at ASC;
