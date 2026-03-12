package server

import (
	"net/http"

	"github.com/perbu/zserve/internal/api"
	"github.com/perbu/zserve/internal/db"
)

func (s *Server) ListGroups(w http.ResponseWriter, r *http.Request, params api.ListGroupsParams) {
	ctx := r.Context()

	limit := int32(50)
	if params.Limit != nil {
		limit = int32(*params.Limit)
	}
	offset := int32(0)
	if params.Offset != nil {
		offset = int32(*params.Offset)
	}

	total, err := s.queries.CountGroups(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rows, err := s.queries.ListGroups(ctx, db.ListGroupsParams{
		QueryLimit:  limit,
		QueryOffset: offset,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	groups := make([]api.Group, 0, len(rows))
	for _, row := range rows {
		groups = append(groups, groupFromRow(row))
	}

	writeJSON(w, http.StatusOK, api.GroupList{
		Groups: groups,
		Total:  total,
		Limit:  int(limit),
		Offset: int(offset),
	})
}

func (s *Server) GetGroup(w http.ResponseWriter, r *http.Request, groupId int64) {
	row, err := s.queries.GetGroup(r.Context(), groupId)
	if err != nil {
		writeError(w, http.StatusNotFound, "group not found")
		return
	}
	writeJSON(w, http.StatusOK, groupFromRow(row))
}

func groupFromRow(row db.ZendeskGroup) api.Group {
	g := api.Group{
		Id:          row.ID,
		Name:        pgtextPtr(row.Name),
		Description: pgtextPtr(row.Description),
		IsPublic:    pgboolPtr(row.IsPublic),
		Default:     pgboolPtr(row.Default),
		Deleted:     pgboolPtr(row.Deleted),
	}
	if row.CreatedAt.Valid {
		t := row.CreatedAt.Time
		g.CreatedAt = &t
	}
	if row.UpdatedAt.Valid {
		t := row.UpdatedAt.Time
		g.UpdatedAt = &t
	}
	return g
}
