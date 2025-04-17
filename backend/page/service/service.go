package service

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/dto/service/input"
	"github.com/trysourcetool/sourcetool/backend/dto/service/output"
	"github.com/trysourcetool/sourcetool/backend/environment"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/page"
	"github.com/trysourcetool/sourcetool/backend/user"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type PageService interface {
	List(context.Context, input.ListPagesInput) (*output.ListPagesOutput, error)
}

type PageServiceCE struct {
	*infra.Dependency
}

func NewPageServiceCE(d *infra.Dependency) *PageServiceCE {
	return &PageServiceCE{Dependency: d}
}

func (s *PageServiceCE) List(ctx context.Context, in input.ListPagesInput) (*output.ListPagesOutput, error) {
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

	users, err := s.Store.User().List(ctx, user.ByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	userGroups, err := s.Store.User().ListGroups(ctx, user.GroupByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	pagesOut := make([]*output.Page, 0, len(pages))
	for _, page := range pages {
		pagesOut = append(pagesOut, output.PageFromModel(page))
	}

	usersOut := make([]*output.User, 0, len(users))
	for _, u := range users {
		usersOut = append(usersOut, output.UserFromModel(u, nil, user.UserOrganizationRoleUnknown))
	}

	userGroupsOut := make([]*output.UserGroup, 0, len(userGroups))
	for _, userGroup := range userGroups {
		userGroupsOut = append(userGroupsOut, output.UserGroupFromModel(userGroup))
	}

	return &output.ListPagesOutput{
		Pages:      pagesOut,
		Groups:     make([]*output.Group, 0),
		GroupPages: make([]*output.GroupPage, 0),
		Users:      usersOut,
		UserGroups: userGroupsOut,
	}, nil
}
