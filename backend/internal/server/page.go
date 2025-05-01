package server

import (
	"context"
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
)

type pageResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Route     string `json:"route"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (s *Server) pageFromModel(page *core.Page) *pageResponse {
	if page == nil {
		return nil
	}

	return &pageResponse{
		ID:        page.ID.String(),
		Name:      page.Name,
		Route:     page.Route,
		CreatedAt: strconv.FormatInt(page.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(page.UpdatedAt.Unix(), 10),
	}
}

type listPagesResponse struct {
	Pages      []*pageResponse      `json:"pages"`
	Groups     []*groupResponse     `json:"groups"`
	GroupPages []*groupPageResponse `json:"groupPages"`
	Users      []*userResponse      `json:"users"`
	UserGroups []*userGroupResponse `json:"userGroups"`
}

func (s *Server) handleListPagesBase(ctx context.Context, env *core.Environment, o *core.Organization) (
	[]*pageResponse,
	[]*userResponse,
	[]*userGroupResponse,
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

	pagesOut := make([]*pageResponse, 0, len(pages))
	for _, page := range pages {
		pagesOut = append(pagesOut, s.pageFromModel(page))
	}

	usersOut := make([]*userResponse, 0, len(users))
	for _, u := range users {
		usersOut = append(usersOut, s.userFromModel(u, core.UserOrganizationRoleUnknown, o))
	}

	userGroupsOut := make([]*userGroupResponse, 0, len(userGroups))
	for _, userGroup := range userGroups {
		userGroupsOut = append(userGroupsOut, s.userGroupFromModel(userGroup))
	}

	return pagesOut, usersOut, userGroupsOut, nil
}
