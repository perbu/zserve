package server

import (
	"net/http"

	"github.com/perbu/zserve/internal/api"
	"github.com/perbu/zserve/internal/db"
)

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request, params api.ListUsersParams) {
	ctx := r.Context()

	limit := int32(50)
	if params.Limit != nil {
		limit = int32(*params.Limit)
	}
	offset := int32(0)
	if params.Offset != nil {
		offset = int32(*params.Offset)
	}

	total, err := s.queries.CountUsers(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rows, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		QueryLimit:  limit,
		QueryOffset: offset,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	users := make([]api.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, userFromRow(row))
	}

	writeJSON(w, http.StatusOK, api.UserList{
		Users:  users,
		Total:  total,
		Limit:  int(limit),
		Offset: int(offset),
	})
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request, userId int64) {
	row, err := s.queries.GetUser(r.Context(), userId)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, userFromRow(row))
}

func userFromRow(row db.ZendeskUser) api.User {
	u := api.User{
		Id:             row.ID,
		Name:           pgtextPtr(row.Name),
		Email:          pgtextPtr(row.Email),
		Role:           pgtextPtr(row.Role),
		OrganizationId: pgint8Ptr(row.OrganizationID),
		Active:         pgboolPtr(row.Active),
		Suspended:      pgboolPtr(row.Suspended),
		TimeZone:       pgtextPtr(row.TimeZone),
		Locale:         pgtextPtr(row.Locale),
		Phone:          pgtextPtr(row.Phone),
	}
	if row.CreatedAt.Valid {
		t := row.CreatedAt.Time
		u.CreatedAt = &t
	}
	if row.UpdatedAt.Valid {
		t := row.UpdatedAt.Time
		u.UpdatedAt = &t
	}
	return u
}
