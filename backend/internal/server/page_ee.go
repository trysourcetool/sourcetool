//go:build ee
// +build ee

package server

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
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

	env, err := s.db.GetEnvironment(ctx, postgres.EnvironmentByID(envID))
	if err != nil {
		return err
	}

	pages, err := s.db.ListPages(ctx, postgres.PageByOrganizationID(o.ID), postgres.PageByEnvironmentID(env.ID), postgres.PageOrderBy(`array_length(p."path", 1), "path"`))
	if err != nil {
		return err
	}

	groups, err := s.db.ListGroups(ctx, postgres.GroupByOrganizationID(o.ID))
	if err != nil {
		return err
	}

	groupPages, err := s.db.ListGroupPages(ctx, postgres.GroupPageByOrganizationID(o.ID), postgres.GroupPageByEnvironmentID(env.ID))
	if err != nil {
		return err
	}

	users, err := s.db.ListUsers(ctx, postgres.UserByOrganizationID(o.ID))
	if err != nil {
		return err
	}

	userGroups, err := s.db.ListUserGroups(ctx, postgres.UserGroupByOrganizationID(o.ID))
	if err != nil {
		return err
	}

	pagesOut := make([]*responses.PageResponse, 0, len(pages))
	for _, page := range pages {
		pagesOut = append(pagesOut, responses.PageFromModel(page))
	}

	groupsOut := make([]*responses.GroupResponse, 0, len(groups))
	for _, group := range groups {
		groupsOut = append(groupsOut, responses.GroupFromModel(group))
	}

	groupPagesOut := make([]*responses.GroupPageResponse, 0, len(groupPages))
	for _, groupPage := range groupPages {
		groupPagesOut = append(groupPagesOut, responses.GroupPageFromModel(groupPage))
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
		Groups:     groupsOut,
		GroupPages: groupPagesOut,
		Users:      usersOut,
		UserGroups: userGroupsOut,
	})
}
