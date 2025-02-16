package page

import (
	"context"
	"strconv"

	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type ServiceCE interface {
	List(context.Context) (*types.ListPagesPayload, error)
}

type ServiceCEImpl struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCEImpl {
	return &ServiceCEImpl{Dependency: d}
}

func (s *ServiceCEImpl) List(ctx context.Context) (*types.ListPagesPayload, error) {
	o := ctxutils.CurrentOrganization(ctx)

	pages, err := s.Store.Page().List(ctx, model.PageByOrganizationID(o.ID), infra.OrderBy(`array_length(p."path", 1), "path"`))
	if err != nil {
		return nil, err
	}

	groups, err := s.Store.Group().List(ctx, model.GroupByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	groupPages, err := s.Store.Group().ListPages(ctx, model.GroupPageByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	users, err := s.Store.User().List(ctx, model.UserByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	userGroups, err := s.Store.User().ListGroups(ctx, model.UserGroupByOrganizationID(o.ID))
	if err != nil {
		return nil, err
	}

	pagesRes := make([]*types.PagePayload, 0, len(pages))
	for _, page := range pages {
		pagesRes = append(pagesRes, &types.PagePayload{
			ID:        page.ID.String(),
			Name:      page.Name,
			Route:     page.Route,
			CreatedAt: strconv.FormatInt(page.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(page.UpdatedAt.Unix(), 10),
		})
	}

	groupsRes := make([]*types.GroupPayload, 0, len(groups))
	for _, group := range groups {
		groupsRes = append(groupsRes, &types.GroupPayload{
			ID:        group.ID.String(),
			Name:      group.Name,
			Slug:      group.Slug,
			CreatedAt: strconv.FormatInt(group.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(group.UpdatedAt.Unix(), 10),
		})
	}

	groupPagesRes := make([]*types.GroupPagePayload, 0, len(groupPages))
	for _, groupPage := range groupPages {
		groupPagesRes = append(groupPagesRes, &types.GroupPagePayload{
			ID:        groupPage.ID.String(),
			GroupID:   groupPage.GroupID.String(),
			PageID:    groupPage.PageID.String(),
			CreatedAt: strconv.FormatInt(groupPage.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(groupPage.UpdatedAt.Unix(), 10),
		})
	}

	usersRes := make([]*types.UserPayload, 0, len(users))
	for _, user := range users {
		usersRes = append(usersRes, &types.UserPayload{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: strconv.FormatInt(user.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(user.UpdatedAt.Unix(), 10),
		})
	}

	userGroupsRes := make([]*types.UserGroupPayload, 0, len(userGroups))
	for _, userGroup := range userGroups {
		userGroupsRes = append(userGroupsRes, &types.UserGroupPayload{
			ID:        userGroup.ID.String(),
			UserID:    userGroup.UserID.String(),
			GroupID:   userGroup.GroupID.String(),
			CreatedAt: strconv.FormatInt(userGroup.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(userGroup.UpdatedAt.Unix(), 10),
		})
	}

	return &types.ListPagesPayload{
		Pages:      pagesRes,
		Groups:     groupsRes,
		GroupPages: groupPagesRes,
		Users:      usersRes,
		UserGroups: userGroupsRes,
	}, nil
}
