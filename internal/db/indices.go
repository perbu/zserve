package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// indices that zserve needs for acceptable read performance.
// These are created with IF NOT EXISTS so they're safe to run on every startup.
var requiredIndices = []string{
	// ticket_comments lookup by ticket
	`CREATE INDEX IF NOT EXISTS idx_ticket_comments_ticket_id ON zendesk.ticket_comments (ticket_id)`,

	// tickets filtered by foreign keys
	`CREATE INDEX IF NOT EXISTS idx_tickets_requester_id ON zendesk.tickets (requester_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tickets_assignee_id ON zendesk.tickets (assignee_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tickets_organization_id ON zendesk.tickets (organization_id)`,

	// tickets sorted by updated_at
	`CREATE INDEX IF NOT EXISTS idx_tickets_updated_at ON zendesk.tickets (updated_at)`,

	// child table joins via _dlt_root_id
	`CREATE INDEX IF NOT EXISTS idx_tickets_tags_dlt_root_id ON zendesk.tickets__tags (_dlt_root_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tickets_collaborator_ids_dlt_root_id ON zendesk.tickets__collaborator_ids (_dlt_root_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tickets_follower_ids_dlt_root_id ON zendesk.tickets__follower_ids (_dlt_root_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tickets_email_cc_ids_dlt_root_id ON zendesk.tickets__email_cc_ids (_dlt_root_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tickets_followup_ids_dlt_root_id ON zendesk.tickets__followup_ids (_dlt_root_id)`,

	// users lookup by id (no PK on real dlt table)
	`CREATE INDEX IF NOT EXISTS idx_users_id ON zendesk.users (id)`,
	`CREATE INDEX IF NOT EXISTS idx_users_organization_id ON zendesk.users (organization_id)`,

	// organizations lookup by id
	`CREATE INDEX IF NOT EXISTS idx_organizations_id ON zendesk.organizations (id)`,

	// groups lookup by id
	`CREATE INDEX IF NOT EXISTS idx_groups_id ON zendesk.groups (id)`,
}

// EnsureIndices creates any missing indices needed by zserve.
func EnsureIndices(ctx context.Context, pool *pgxpool.Pool) error {
	for _, ddl := range requiredIndices {
		if _, err := pool.Exec(ctx, ddl); err != nil {
			return fmt.Errorf("creating index: %w\n  DDL: %s", err, ddl)
		}
	}
	slog.Info("database indices verified", "count", len(requiredIndices))
	return nil
}
