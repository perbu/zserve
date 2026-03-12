package server

import (
	"context"
	"log/slog"

	"github.com/perbu/zserve/internal/api"
)

// enrichTickets resolves requester_name, assignee_name, submitter_name,
// organization_name, and group_name for a slice of tickets via batch lookups.
func (s *Server) enrichTickets(ctx context.Context, tickets []api.Ticket) {
	if len(tickets) == 0 {
		return
	}

	// Collect unique IDs
	userIDs := make(map[int64]bool)
	orgIDs := make(map[int64]bool)
	groupIDs := make(map[int64]bool)

	for _, t := range tickets {
		if t.RequesterId != nil {
			userIDs[*t.RequesterId] = true
		}
		if t.AssigneeId != nil {
			userIDs[*t.AssigneeId] = true
		}
		if t.SubmitterId != nil {
			userIDs[*t.SubmitterId] = true
		}
		if t.OrganizationId != nil {
			orgIDs[*t.OrganizationId] = true
		}
		if t.GroupId != nil {
			groupIDs[*t.GroupId] = true
		}
	}

	// Batch fetch names
	userNames := s.fetchUserNames(ctx, userIDs)
	orgNames := s.fetchOrgNames(ctx, orgIDs)
	grpNames := s.fetchGroupNames(ctx, groupIDs)

	// Assign names
	for i := range tickets {
		if tickets[i].RequesterId != nil {
			if name, ok := userNames[*tickets[i].RequesterId]; ok {
				tickets[i].RequesterName = &name
			}
		}
		if tickets[i].AssigneeId != nil {
			if name, ok := userNames[*tickets[i].AssigneeId]; ok {
				tickets[i].AssigneeName = &name
			}
		}
		if tickets[i].SubmitterId != nil {
			if name, ok := userNames[*tickets[i].SubmitterId]; ok {
				tickets[i].SubmitterName = &name
			}
		}
		if tickets[i].OrganizationId != nil {
			if name, ok := orgNames[*tickets[i].OrganizationId]; ok {
				tickets[i].OrganizationName = &name
			}
		}
		if tickets[i].GroupId != nil {
			if name, ok := grpNames[*tickets[i].GroupId]; ok {
				tickets[i].GroupName = &name
			}
		}
	}
}

// enrichComments resolves author_name for a slice of comments via batch lookup.
func (s *Server) enrichComments(ctx context.Context, comments []api.Comment) {
	if len(comments) == 0 {
		return
	}

	userIDs := make(map[int64]bool)
	for _, c := range comments {
		if c.AuthorId != nil {
			userIDs[*c.AuthorId] = true
		}
	}

	userNames := s.fetchUserNames(ctx, userIDs)

	for i := range comments {
		if comments[i].AuthorId != nil {
			if name, ok := userNames[*comments[i].AuthorId]; ok {
				comments[i].AuthorName = &name
			}
		}
	}
}

func (s *Server) fetchUserNames(ctx context.Context, idSet map[int64]bool) map[int64]string {
	result := make(map[int64]string)
	if len(idSet) == 0 {
		return result
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	rows, err := s.queries.GetUserNames(ctx, ids)
	if err != nil {
		slog.Warn("failed to fetch user names for enrichment", "error", err)
		return result
	}
	for _, r := range rows {
		if r.Name.Valid {
			result[r.ID] = r.Name.String
		}
	}
	return result
}

func (s *Server) fetchOrgNames(ctx context.Context, idSet map[int64]bool) map[int64]string {
	result := make(map[int64]string)
	if len(idSet) == 0 {
		return result
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	rows, err := s.queries.GetOrganizationNames(ctx, ids)
	if err != nil {
		slog.Warn("failed to fetch organization names for enrichment", "error", err)
		return result
	}
	for _, r := range rows {
		if r.Name.Valid {
			result[r.ID] = r.Name.String
		}
	}
	return result
}

func (s *Server) fetchGroupNames(ctx context.Context, idSet map[int64]bool) map[int64]string {
	result := make(map[int64]string)
	if len(idSet) == 0 {
		return result
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	rows, err := s.queries.GetGroupNames(ctx, ids)
	if err != nil {
		slog.Warn("failed to fetch group names for enrichment", "error", err)
		return result
	}
	for _, r := range rows {
		if r.Name.Valid {
			result[r.ID] = r.Name.String
		}
	}
	return result
}
