package group

import (
	"context"
	"errors"
	"strconv"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/authz"
	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/server/http/types"
)

type ServiceCE interface {
	Get(context.Context, types.GetGroupInput) (*types.GetGroupPayload, error)
	List(context.Context) (*types.ListGroupsPayload, error)
	Create(context.Context, types.CreateGroupInput) (*types.CreateGroupPayload, error)
	Update(context.Context, types.UpdateGroupInput) (*types.UpdateGroupPayload, error)
	Delete(context.Context, types.DeleteGroupInput) (*types.DeleteGroupPayload, error)
}

type ServiceCEImpl struct {
	*infra.Dependency
}

func NewServiceCE(d *infra.Dependency) *ServiceCEImpl {
	return &ServiceCEImpl{Dependency: d}
}

func (s *ServiceCEImpl) Get(ctx context.Context, in types.GetGroupInput) (*types.GetGroupPayload, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	group, err := s.Store.Group().Get(ctx, model.GroupByOrganizationID(currentOrg.ID), model.GroupByID(groupID))
	if err != nil {
		return nil, err
	}

	return &types.GetGroupPayload{
		Group: &types.GroupPayload{
			ID:        group.ID.String(),
			Name:      group.Name,
			Slug:      group.Slug,
			CreatedAt: strconv.FormatInt(group.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(group.UpdatedAt.Unix(), 10),
		},
	}, nil
}

func (s *ServiceCEImpl) List(ctx context.Context) (*types.ListGroupsPayload, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	groups, err := s.Store.Group().List(ctx, model.GroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	users, err := s.Store.User().List(ctx, model.UserByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userGroups, err := s.Store.User().ListGroups(ctx, model.UserGroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
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

	return &types.ListGroupsPayload{
		Groups:     groupsRes,
		Users:      usersRes,
		UserGroups: userGroupsRes,
	}, nil
}

func (s *ServiceCEImpl) Create(ctx context.Context, in types.CreateGroupInput) (*types.CreateGroupPayload, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditGroup); err != nil {
		return nil, err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)

	slugExists, err := s.Store.Group().IsSlugExistsInOrganization(ctx, currentOrg.ID, in.Slug)
	if err != nil {
		return nil, err
	}
	if slugExists {
		return nil, errdefs.ErrGroupSlugAlreadyExists(errors.New("slug already exists"))
	}

	if !validateSlug(in.Slug) {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid slug"))
	}

	group := &model.Group{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: currentOrg.ID,
		Name:           in.Name,
		Slug:           in.Slug,
	}

	userIDs := make([]uuid.UUID, 0, len(in.UserIDs))
	for _, userID := range in.UserIDs {
		id, err := uuid.FromString(userID)
		if err != nil {
			return nil, errdefs.ErrInvalidArgument(err)
		}
		userIDs = append(userIDs, id)
	}

	userGroups := make([]*model.UserGroup, 0, len(userIDs))
	for _, userID := range userIDs {
		userGroups = append(userGroups, &model.UserGroup{
			ID:      uuid.Must(uuid.NewV4()),
			UserID:  userID,
			GroupID: group.ID,
		})
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.Group().Create(ctx, group); err != nil {
			return err
		}

		if err := tx.User().BulkInsertGroups(ctx, userGroups); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	group, err = s.Store.Group().Get(ctx, model.GroupByID(group.ID))
	if err != nil {
		return nil, err
	}

	return &types.CreateGroupPayload{
		Group: &types.GroupPayload{
			ID:        group.ID.String(),
			Name:      group.Name,
			Slug:      group.Slug,
			CreatedAt: strconv.FormatInt(group.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(group.UpdatedAt.Unix(), 10),
		},
	}, nil
}

func (s *ServiceCEImpl) Update(ctx context.Context, in types.UpdateGroupInput) (*types.UpdateGroupPayload, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditGroup); err != nil {
		return nil, err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	group, err := s.Store.Group().Get(ctx, model.GroupByID(groupID), model.GroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		group.Name = conv.SafeValue(in.Name)
	}

	userIDs := make([]uuid.UUID, 0, len(in.UserIDs))
	for _, userID := range in.UserIDs {
		id, err := uuid.FromString(userID)
		if err != nil {
			return nil, errdefs.ErrInvalidArgument(err)
		}
		userIDs = append(userIDs, id)
	}

	userGroups := make([]*model.UserGroup, 0, len(userIDs))
	for _, userID := range userIDs {
		userGroups = append(userGroups, &model.UserGroup{
			ID:      uuid.Must(uuid.NewV4()),
			UserID:  userID,
			GroupID: group.ID,
		})
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.Group().Update(ctx, group); err != nil {
			return err
		}

		existingGroups, err := tx.User().ListGroups(ctx, model.UserGroupByGroupID(group.ID))
		if err != nil {
			return err
		}

		if err := tx.User().BulkDeleteGroups(ctx, existingGroups); err != nil {
			return err
		}

		if err := tx.User().BulkInsertGroups(ctx, userGroups); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	group, err = s.Store.Group().Get(ctx, model.GroupByID(group.ID))
	if err != nil {
		return nil, err
	}

	return &types.UpdateGroupPayload{
		Group: &types.GroupPayload{
			ID:        group.ID.String(),
			Name:      group.Name,
			Slug:      group.Slug,
			CreatedAt: strconv.FormatInt(group.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(group.UpdatedAt.Unix(), 10),
		},
	}, nil
}

func (s *ServiceCEImpl) Delete(ctx context.Context, in types.DeleteGroupInput) (*types.DeleteGroupPayload, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditGroup); err != nil {
		return nil, err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	group, err := s.Store.Group().Get(ctx, model.GroupByID(groupID), model.GroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if err := s.Store.RunTransaction(func(tx infra.Transaction) error {
		if err := tx.Group().Delete(ctx, group); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &types.DeleteGroupPayload{
		Group: &types.GroupPayload{
			ID:        group.ID.String(),
			Name:      group.Name,
			Slug:      group.Slug,
			CreatedAt: strconv.FormatInt(group.CreatedAt.Unix(), 10),
			UpdatedAt: strconv.FormatInt(group.UpdatedAt.Unix(), 10),
		},
	}, nil
}
