package page

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
	"github.com/trysourcetool/sourcetool/backend/utils/ctxutil"
)

type Service interface {
	List(context.Context, *dto.ListPagesInput) (*dto.ListPagesOutput, error)
}

type ServiceCE struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCE {
	return &ServiceCE{Dependency: d}
}

func (s *ServiceCE) List(ctx context.Context, in *dto.ListPagesInput) (*dto.ListPagesOutput, error) {
	o := ctxutil.CurrentOrganization(ctx)

	envID, err := uuid.FromString(in.EnvironmentID)
	if err != nil {
		return nil, err
	}

	env, err := s.Store.Environment().Get(ctx, storeopts.EnvironmentByID(envID))
	if err != nil {
		return nil, err
	}

	pages, err := s.Store.Page().List(ctx, storeopts.PageByOrganizationID(o.ID), storeopts.PageByEnvironmentID(env.ID), storeopts.PageOrderBy(`array_length(p."path", 1), "path"`))
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
		Groups:     make([]*dto.Group, 0),
		GroupPages: make([]*dto.GroupPage, 0),
		Users:      usersOut,
		UserGroups: userGroupsOut,
	}, nil
}
