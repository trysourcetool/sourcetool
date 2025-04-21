package server

import (
	"context"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/server/responses"
)

func (s *Server) listPagesBase(ctx context.Context, env *core.Environment, o *core.Organization) (
	[]*responses.PageResponse,
	[]*responses.UserResponse,
	[]*responses.UserGroupResponse,
	error,
) {
	pages, err := s.db.Page().List(ctx, database.PageByOrganizationID(o.ID), database.PageByEnvironmentID(env.ID), database.PageOrderBy(`array_length(p."path", 1), "path"`))
	if err != nil {
		return nil, nil, nil, err
	}

	users, err := s.db.User().List(ctx, database.UserByOrganizationID(o.ID))
	if err != nil {
		return nil, nil, nil, err
	}

	userGroups, err := s.db.User().ListGroups(ctx, database.UserGroupByOrganizationID(o.ID))
	if err != nil {
		return nil, nil, nil, err
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

	return pagesOut, usersOut, userGroupsOut, nil
}
