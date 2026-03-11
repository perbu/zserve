package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/perbu/zserve/internal/api"
	"github.com/perbu/zserve/internal/db"
)

type Server struct {
	queries *db.Queries
}

func New(queries *db.Queries) *Server {
	return &Server{queries: queries}
}

func (s *Server) ListTickets(w http.ResponseWriter, r *http.Request, params api.ListTicketsParams) {
	ctx := r.Context()

	limit := int32(50)
	if params.Limit != nil {
		limit = int32(*params.Limit)
	}
	offset := int32(0)
	if params.Offset != nil {
		offset = int32(*params.Offset)
	}

	f := buildFilterParams(params)

	countParams := db.CountTicketsParams{
		Status: f.Status, TicketType: f.TicketType, Priority: f.Priority,
		RequesterID: f.RequesterID, AssigneeID: f.AssigneeID, SubmitterID: f.SubmitterID,
		OrganizationID: f.OrganizationID, GroupID: f.GroupID, BrandID: f.BrandID,
		Search: f.Search, CreatedAfter: f.CreatedAfter, CreatedBefore: f.CreatedBefore,
		Tag: f.Tag,
	}
	total, err := s.queries.CountTickets(ctx, countParams)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sortBy := "created_at"
	if params.SortBy != nil {
		sortBy = string(*params.SortBy)
	}
	sortOrder := "desc"
	if params.SortOrder != nil {
		sortOrder = string(*params.SortOrder)
	}

	rows, err := s.listTickets(ctx, sortBy, sortOrder, limit, offset, f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tickets := make([]api.Ticket, 0, len(rows))
	for _, row := range rows {
		tickets = append(tickets, ticketFromRow(row))
	}

	writeJSON(w, http.StatusOK, api.TicketList{
		Tickets: tickets,
		Total:   total,
		Limit:   int(limit),
		Offset:  int(offset),
	})
}

func (s *Server) GetTicket(w http.ResponseWriter, r *http.Request, ticketId int64) {
	ctx := r.Context()

	row, err := s.queries.GetTicket(ctx, ticketId)
	if err != nil {
		writeError(w, http.StatusNotFound, "ticket not found")
		return
	}

	ticket := ticketFromRow(row)
	detail := api.TicketDetail{}
	data, _ := json.Marshal(ticket)
	json.Unmarshal(data, &detail)

	tags := s.fetchTags(ctx, ticketId)
	collabs := s.fetchCollaborators(ctx, ticketId)
	followers := s.fetchFollowers(ctx, ticketId)
	emailCCs := s.fetchEmailCCs(ctx, ticketId)
	followups := s.fetchFollowups(ctx, ticketId)

	detail.TagList = &tags
	detail.CollaboratorIds = &collabs
	detail.FollowerIds = &followers
	detail.EmailCcIds = &emailCCs
	detail.FollowupIds = &followups

	writeJSON(w, http.StatusOK, detail)
}

func (s *Server) GetTicketTags(w http.ResponseWriter, r *http.Request, ticketId int64) {
	tags := s.fetchTags(r.Context(), ticketId)
	writeJSON(w, http.StatusOK, tags)
}

func (s *Server) GetTicketCollaborators(w http.ResponseWriter, r *http.Request, ticketId int64) {
	writeJSON(w, http.StatusOK, s.fetchCollaborators(r.Context(), ticketId))
}

func (s *Server) GetTicketFollowers(w http.ResponseWriter, r *http.Request, ticketId int64) {
	writeJSON(w, http.StatusOK, s.fetchFollowers(r.Context(), ticketId))
}

func (s *Server) GetTicketEmailCCs(w http.ResponseWriter, r *http.Request, ticketId int64) {
	writeJSON(w, http.StatusOK, s.fetchEmailCCs(r.Context(), ticketId))
}

func (s *Server) GetTicketFollowups(w http.ResponseWriter, r *http.Request, ticketId int64) {
	writeJSON(w, http.StatusOK, s.fetchFollowups(r.Context(), ticketId))
}

// helpers

func (s *Server) fetchTags(ctx context.Context, ticketId int64) []string {
	rows, err := s.queries.GetTicketTags(ctx, ticketId)
	if err != nil {
		return []string{}
	}
	tags := make([]string, 0, len(rows))
	for _, r := range rows {
		if r.Valid {
			tags = append(tags, r.String)
		}
	}
	return tags
}

func (s *Server) fetchCollaborators(ctx context.Context, ticketId int64) []int64 {
	rows, err := s.queries.GetTicketCollaborators(ctx, ticketId)
	if err != nil {
		return []int64{}
	}
	return pgIntsToSlice(rows)
}

func (s *Server) fetchFollowers(ctx context.Context, ticketId int64) []int64 {
	rows, err := s.queries.GetTicketFollowers(ctx, ticketId)
	if err != nil {
		return []int64{}
	}
	return pgIntsToSlice(rows)
}

func (s *Server) fetchEmailCCs(ctx context.Context, ticketId int64) []int64 {
	rows, err := s.queries.GetTicketEmailCCs(ctx, ticketId)
	if err != nil {
		return []int64{}
	}
	return pgIntsToSlice(rows)
}

func (s *Server) fetchFollowups(ctx context.Context, ticketId int64) []int64 {
	rows, err := s.queries.GetTicketFollowups(ctx, ticketId)
	if err != nil {
		return []int64{}
	}
	return pgIntsToSlice(rows)
}

func pgIntsToSlice(rows []pgtype.Int8) []int64 {
	out := make([]int64, 0, len(rows))
	for _, r := range rows {
		if r.Valid {
			out = append(out, r.Int64)
		}
	}
	return out
}

type filterFields struct {
	Status         pgtype.Text
	TicketType     pgtype.Text
	Priority       pgtype.Text
	RequesterID    pgtype.Int8
	AssigneeID     pgtype.Int8
	SubmitterID    pgtype.Int8
	OrganizationID pgtype.Int8
	GroupID        pgtype.Int8
	BrandID        pgtype.Int8
	Search         pgtype.Text
	CreatedAfter   pgtype.Timestamptz
	CreatedBefore  pgtype.Timestamptz
	Tag            pgtype.Text
}

func buildFilterParams(params api.ListTicketsParams) filterFields {
	var f filterFields
	if params.Status != nil {
		f.Status = pgtype.Text{String: string(*params.Status), Valid: true}
	}
	if params.Type != nil {
		f.TicketType = pgtype.Text{String: string(*params.Type), Valid: true}
	}
	if params.Priority != nil {
		f.Priority = pgtype.Text{String: string(*params.Priority), Valid: true}
	}
	if params.RequesterId != nil {
		f.RequesterID = pgtype.Int8{Int64: *params.RequesterId, Valid: true}
	}
	if params.AssigneeId != nil {
		f.AssigneeID = pgtype.Int8{Int64: *params.AssigneeId, Valid: true}
	}
	if params.SubmitterId != nil {
		f.SubmitterID = pgtype.Int8{Int64: *params.SubmitterId, Valid: true}
	}
	if params.OrganizationId != nil {
		f.OrganizationID = pgtype.Int8{Int64: *params.OrganizationId, Valid: true}
	}
	if params.GroupId != nil {
		f.GroupID = pgtype.Int8{Int64: *params.GroupId, Valid: true}
	}
	if params.BrandId != nil {
		f.BrandID = pgtype.Int8{Int64: *params.BrandId, Valid: true}
	}
	if params.Search != nil {
		f.Search = pgtype.Text{String: *params.Search, Valid: true}
	}
	if params.CreatedAfter != nil {
		f.CreatedAfter = pgtype.Timestamptz{Time: *params.CreatedAfter, Valid: true}
	}
	if params.CreatedBefore != nil {
		f.CreatedBefore = pgtype.Timestamptz{Time: *params.CreatedBefore, Valid: true}
	}
	if params.Tag != nil {
		f.Tag = pgtype.Text{String: *params.Tag, Valid: true}
	}
	return f
}

func (s *Server) listTickets(ctx context.Context, sortBy, sortOrder string, limit, offset int32, f filterFields) ([]db.ZendeskTicket, error) {
	switch sortBy + "_" + sortOrder {
	case "updated_at_asc":
		return s.queries.ListTicketsByUpdatedAtAsc(ctx, db.ListTicketsByUpdatedAtAscParams{
			Status: f.Status, TicketType: f.TicketType, Priority: f.Priority,
			RequesterID: f.RequesterID, AssigneeID: f.AssigneeID, SubmitterID: f.SubmitterID,
			OrganizationID: f.OrganizationID, GroupID: f.GroupID, BrandID: f.BrandID,
			Search: f.Search, CreatedAfter: f.CreatedAfter, CreatedBefore: f.CreatedBefore,
			Tag: f.Tag, QueryLimit: limit, QueryOffset: offset,
		})
	case "updated_at_desc":
		return s.queries.ListTicketsByUpdatedAtDesc(ctx, db.ListTicketsByUpdatedAtDescParams{
			Status: f.Status, TicketType: f.TicketType, Priority: f.Priority,
			RequesterID: f.RequesterID, AssigneeID: f.AssigneeID, SubmitterID: f.SubmitterID,
			OrganizationID: f.OrganizationID, GroupID: f.GroupID, BrandID: f.BrandID,
			Search: f.Search, CreatedAfter: f.CreatedAfter, CreatedBefore: f.CreatedBefore,
			Tag: f.Tag, QueryLimit: limit, QueryOffset: offset,
		})
	case "created_at_asc":
		return s.queries.ListTicketsByCreatedAtAsc(ctx, db.ListTicketsByCreatedAtAscParams{
			Status: f.Status, TicketType: f.TicketType, Priority: f.Priority,
			RequesterID: f.RequesterID, AssigneeID: f.AssigneeID, SubmitterID: f.SubmitterID,
			OrganizationID: f.OrganizationID, GroupID: f.GroupID, BrandID: f.BrandID,
			Search: f.Search, CreatedAfter: f.CreatedAfter, CreatedBefore: f.CreatedBefore,
			Tag: f.Tag, QueryLimit: limit, QueryOffset: offset,
		})
	default: // created_at_desc
		return s.queries.ListTicketsByCreatedAtDesc(ctx, db.ListTicketsByCreatedAtDescParams{
			Status: f.Status, TicketType: f.TicketType, Priority: f.Priority,
			RequesterID: f.RequesterID, AssigneeID: f.AssigneeID, SubmitterID: f.SubmitterID,
			OrganizationID: f.OrganizationID, GroupID: f.GroupID, BrandID: f.BrandID,
			Search: f.Search, CreatedAfter: f.CreatedAfter, CreatedBefore: f.CreatedBefore,
			Tag: f.Tag, QueryLimit: limit, QueryOffset: offset,
		})
	}
}

func pgtextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func pgint8Ptr(i pgtype.Int8) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

func pgboolPtr(b pgtype.Bool) *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

func ticketFromRow(row db.ZendeskTicket) api.Ticket {
	t := api.Ticket{
		Id:                   row.ID,
		Subject:              pgtextPtr(row.Subject),
		RawSubject:           pgtextPtr(row.RawSubject),
		Description:          pgtextPtr(row.Description),
		Status:               pgtextPtr(row.Status),
		Type:                 pgtextPtr(row.Type),
		Priority:             pgtextPtr(row.Priority),
		Recipient:            pgtextPtr(row.Recipient),
		RequesterId:          pgint8Ptr(row.RequesterID),
		SubmitterId:          pgint8Ptr(row.SubmitterID),
		AssigneeId:           pgint8Ptr(row.AssigneeID),
		OrganizationId:       pgint8Ptr(row.OrganizationID),
		GroupId:              pgint8Ptr(row.GroupID),
		BrandId:              pgint8Ptr(row.BrandID),
		TicketFormId:         pgint8Ptr(row.TicketFormID),
		CustomStatusId:       pgint8Ptr(row.CustomStatusID),
		ExternalId:           pgtextPtr(row.ExternalID),
		Url:                  pgtextPtr(row.Url),
		HasIncidents:         pgboolPtr(row.HasIncidents),
		IsPublic:             pgboolPtr(row.IsPublic),
		AllowChannelback:     pgboolPtr(row.AllowChannelback),
		AllowAttachments:     pgboolPtr(row.AllowAttachments),
		FromMessagingChannel: pgboolPtr(row.FromMessagingChannel),
		Tags:                 pgtextPtr(row.Tags),
	}

	// Timestamps
	if row.CreatedAt.Valid {
		ca := row.CreatedAt.Time
		t.CreatedAt = &ca
	}
	if row.UpdatedAt.Valid {
		ua := row.UpdatedAt.Time
		t.UpdatedAt = &ua
	}

	// Custom fields (JSON)
	if len(row.CustomFields) > 0 {
		var cf interface{}
		if json.Unmarshal(row.CustomFields, &cf) == nil {
			t.CustomFields = cf
		}
	}

	// Satisfaction rating
	sr := api.SatisfactionRating{
		Score:    pgtextPtr(row.SatisfactionRatingScore),
		Id:       pgint8Ptr(row.SatisfactionRatingID),
		Reason:   pgtextPtr(row.SatisfactionRatingReason),
		ReasonId: pgint8Ptr(row.SatisfactionRatingReasonID),
		Comment:  pgtextPtr(row.SatisfactionRatingComment),
	}
	if sr.Score != nil || sr.Id != nil || sr.Reason != nil || sr.ReasonId != nil || sr.Comment != nil {
		t.SatisfactionRating = &sr
	}

	// Via
	v := api.Via{
		Channel:   pgtextPtr(row.ViaChannel),
		SourceRel: pgtextPtr(row.ViaSourceRel),
	}
	hasVia := v.Channel != nil || v.SourceRel != nil
	if len(row.ViaSourceFrom) > 0 {
		var sf interface{}
		if json.Unmarshal(row.ViaSourceFrom, &sf) == nil {
			v.SourceFrom = &sf
			hasVia = true
		}
	}
	if len(row.ViaSourceTo) > 0 {
		var st interface{}
		if json.Unmarshal(row.ViaSourceTo, &st) == nil {
			v.SourceTo = &st
			hasVia = true
		}
	}
	if hasVia {
		t.Via = &v
	}

	return t
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, api.Error{Message: msg})
}

// Middleware returns an http.Handler that optionally checks the API key.
func Middleware(apiKey string, next http.Handler) http.Handler {
	if apiKey == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != apiKey {
			writeError(w, http.StatusUnauthorized, "invalid or missing API key")
			return
		}
		next.ServeHTTP(w, r)
	})
}
