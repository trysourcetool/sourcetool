package page

import (
	"context"

	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

type serviceEE struct {
	*infra.Dependency
	*page.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: page.NewServiceCE(
			infra.NewDependency(d.Store, d.Mailer),
		),
	}
}

func (s *serviceEE) List(ctx context.Context) (*dto.ListPagesOutput, error) {
	o := ctxutils.CurrentOrganization(ctx)

	pages, err := s.Store.Page().List(ctx, storeopts.PageByOrganizationID(o.ID), storeopts.PageOrderBy(`array_length(p."path", 1), "path"`))
	if err != nil {
		return nil, err
	}

	groups, err := s.Store.Group().List(ctx, storeopts.GroupByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	groupPages, err := s.Store.Group().ListPages(ctx, storeopts.GroupPageByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	users, err := s.Store.User().List(ctx, storeopts.UserByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	userGroups, err := s.Store.User().ListGroups(ctx, storeopts.UserGroupByOrganizationID(o.ID))
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
	for _, user := range users {
		usersOut = append(usersOut, dto.UserFromModel(user, nil, model.UserOrganizationRoleUnknown))
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
