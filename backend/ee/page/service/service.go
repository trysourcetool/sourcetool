package service

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/page/service"
	"github.com/trysourcetool/sourcetool/backend/user"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type pageServiceEE struct {
	*infra.Dependency
	*service.PageServiceCE
}

func NewPageServiceEE(d *infra.Dependency) *pageServiceEE {
	return &pageServiceEE{
		Dependency: d,
		PageServiceCE: service.NewPageServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}

func (s *pageServiceEE) List(ctx context.Context, in *dto.ListPagesInput) (*dto.ListPagesOutput, error) {
	o := ctxutil.CurrentOrganization(ctx)

	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, err
	}

	env, err := s.Store.Environment().Get(ctx, environment.ByID(envID))
	if err != nil {
		return nil, err
	}

	pages, err := s.Store.Page().List(ctx, page.ByOrganizationID(o.ID), page.ByEnvironmentID(env.ID), page.OrderBy(`array_length(p."path", 1), "path"`))
	if err != nil {
		return nil, err
	}

	groups, err := s.Store.Group().List(ctx, group.ByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	groupPages, err := s.Store.Group().ListPages(ctx, group.PageByOrganizationID(o.ID), group.PageByEnvironmentID(env.ID))
	if err != nil {
		return nil, err
	}

	users, err := s.Store.User().List(ctx, user.ByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	userGroups, err := s.Store.User().ListGroups(ctx, user.GroupByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	pagesOut := make([]*dto.Page, 0, len(pages))
	for _, page := range pages {
		pagesOut = append(pagesOut, dto.PageFromModel(page))
	}

	groupsOut := make([]*dto.Group, 0, len(groups))
	for _, group := range groups {
		groupsOut = append(groupsOut, dto.GroupFromModel(group))
	}

	groupPagesOut := make([]*dto.GroupPage, 0, len(groupPages))
	for _, groupPage := range groupPages {
		groupPagesOut = append(groupPagesOut, dto.GroupPageFromModel(groupPage))
	}

	usersOut := make([]*dto.User, 0, len(users))
	for _, u := range users {
		usersOut = append(usersOut, dto.UserFromModel(u, nil, user.UserOrganizationRoleUnknown))
	}

	userGroupsOut := make([]*dto.UserGroup, 0, len(userGroups))
	for _, userGroup := range userGroups {
		userGroupsOut = append(userGroupsOut, dto.UserGroupFromModel(userGroup))
	}

	return &dto.ListPagesOutput{
		Pages:      pagesOut,
		Groups:     groupsOut,
		GroupPages: groupPagesOut,
		Users:      usersOut,
		UserGroups: userGroupsOut,
	}, nil
}
