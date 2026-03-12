package server

import (
	"net/http"

	"github.com/perbu/zserve/internal/api"
	"github.com/perbu/zserve/internal/db"
)

func (s *Server) ListOrganizations(w http.ResponseWriter, r *http.Request, params api.ListOrganizationsParams) {
	ctx := r.Context()

	limit := int32(50)
	if params.Limit != nil {
		limit = int32(*params.Limit)
	}
	offset := int32(0)
	if params.Offset != nil {
		offset = int32(*params.Offset)
	}

	total, err := s.queries.CountOrganizations(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rows, err := s.queries.ListOrganizations(ctx, db.ListOrganizationsParams{
		QueryLimit:  limit,
		QueryOffset: offset,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	orgs := make([]api.Organization, 0, len(rows))
	for _, row := range rows {
		orgs = append(orgs, organizationFromRow(row))
	}

	writeJSON(w, http.StatusOK, api.OrganizationList{
		Organizations: orgs,
		Total:         total,
		Limit:         int(limit),
		Offset:        int(offset),
	})
}

func (s *Server) GetOrganization(w http.ResponseWriter, r *http.Request, organizationId int64) {
	row, err := s.queries.GetOrganization(r.Context(), organizationId)
	if err != nil {
		writeError(w, http.StatusNotFound, "organization not found")
		return
	}
	writeJSON(w, http.StatusOK, organizationFromRow(row))
}

func organizationFromRow(row db.ZendeskOrganization) api.Organization {
	o := api.Organization{
		Id:             row.ID,
		Name:           pgtextPtr(row.Name),
		SharedTickets:  pgboolPtr(row.SharedTickets),
		SharedComments: pgboolPtr(row.SharedComments),
		Details:        pgtextPtr(row.Details),
		Notes:          pgtextPtr(row.Notes),
		GroupId:        pgint8Ptr(row.GroupID),
	}
	if row.CreatedAt.Valid {
		t := row.CreatedAt.Time
		o.CreatedAt = &t
	}
	if row.UpdatedAt.Valid {
		t := row.UpdatedAt.Time
		o.UpdatedAt = &t
	}
	return o
}
