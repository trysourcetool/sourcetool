//go:build !ee
// +build !ee

package server

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) listPages(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	envIDReq := r.URL.Query().Get("environmentId")
	if envIDReq == "" {
		return errdefs.ErrInvalidArgument(errors.New("environmentId is required"))
	}

	o := internal.CurrentOrganization(ctx)

	envID, err := uuid.FromString(envIDReq)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	env, err := s.db.Environment().Get(ctx, database.EnvironmentByID(envID))
	if err != nil {
		return err
	}

	pages, err := s.db.Page().List(ctx, database.PageByOrganizationID(o.ID), database.PageByEnvironmentID(env.ID), database.PageOrderBy(`array_length(p."path", 1), "path"`))
	if err != nil {
		return err
	}

	users, err := s.db.User().List(ctx, database.UserByOrganizationID(o.ID))
	if err != nil {
		return err
	}

	userGroups, err := s.db.User().ListGroups(ctx, database.UserGroupByOrganizationID(o.ID))
	if err != nil {
		return err
	}

	pagesOut := make([]*responses.PageResponse, 0, len(pages))
	for _, page := range pages {
		pagesOut = append(pagesOut, responses.PageFromModel(page))
	}

	usersOut := make([]*responses.UserResponse, 0, len(users))
	for _, u := range users {
		usersOut = append(usersOut, responses.UserFromModel(u, core.UserOrganizationRoleUnknown, o))
	}

	userGroupsOut := make([]*responses.UserGroupResponse, 0, len(userGroups))
	for _, userGroup := range userGroups {
		userGroupsOut = append(userGroupsOut, responses.UserGroupFromModel(userGroup))
	}

	return s.renderJSON(w, http.StatusOK, responses.ListPagesResponse{
		Pages:      pagesOut,
		Groups:     make([]*responses.GroupResponse, 0),
		GroupPages: make([]*responses.GroupPageResponse, 0),
		Users:      usersOut,
		UserGroups: userGroupsOut,
	})
}
