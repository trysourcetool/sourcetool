package page

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/environment"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/page"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type Service interface {
	List(context.Context, dto.ListPagesInput) (*dto.ListPagesOutput, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) List(ctx context.Context, in dto.ListPagesInput) (*dto.ListPagesOutput, error) {
	o := ctxutil.CurrentOrganization(ctx)

	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, err
	}

	env, err := s.Repository.Environment().Get(ctx, environment.ByID(envID))
	if err != nil {
		return nil, err
	}

	pages, err := s.Repository.Page().List(ctx, page.ByOrganizationID(o.ID), page.ByEnvironmentID(env.ID), page.OrderBy(`array_length(p."path", 1), "path"`))
	if err != nil {
		return nil, err
	}

	users, err := s.Repository.User().List(ctx, user.ByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	userGroups, err := s.Repository.User().ListGroups(ctx, user.GroupByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	pagesOut := make([]*dto.Page, 0, len(pages))
	for _, page := range pages {
		pagesOut = append(pagesOut, dto.PageFromModel(page))
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
		Groups:     make([]*dto.Group, 0),
		GroupPages: make([]*dto.GroupPage, 0),
		Users:      usersOut,
		UserGroups: userGroupsOut,
	}, nil
}
