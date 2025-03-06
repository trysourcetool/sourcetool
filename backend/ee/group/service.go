package group

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/authz"
	"github.com/trysourcetool/sourcetool/backend/conv"
	"github.com/trysourcetool/sourcetool/backend/ctxutils"
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/group"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

type serviceEE struct {
	*infra.Dependency
	*group.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: group.NewServiceCE(
			infra.NewDependency(d.Store, d.Signer, d.Mailer),
		),
	}
}

func (s *serviceEE) Get(ctx context.Context, in dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	group, err := s.Store.Group().Get(ctx, storeopts.GroupByOrganizationID(currentOrg.ID), storeopts.GroupByID(groupID))
	if err != nil {
		return nil, err
	}

	return &dto.GetGroupOutput{
		Group: dto.GroupFromModel(group),
	}, nil
}

func (s *serviceEE) List(ctx context.Context) (*dto.ListGroupsOutput, error) {
	currentOrg := ctxutils.CurrentOrganization(ctx)
	groups, err := s.Store.Group().List(ctx, storeopts.GroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	users, err := s.Store.User().List(ctx, storeopts.UserByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userGroups, err := s.Store.User().ListGroups(ctx, storeopts.UserGroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userIDs := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	orgAccesses, err := s.Store.User().ListOrganizationAccesses(ctx, storeopts.UserOrganizationAccessByUserIDs(userIDs))
	if err != nil {
		return nil, err
	}

	orgAccessesMap := make(map[uuid.UUID]*model.UserOrganizationAccess)
	for _, orgAccess := range orgAccesses {
		orgAccessesMap[orgAccess.UserID] = orgAccess
	}

	groupsOut := make([]*dto.Group, 0, len(groups))
	for _, group := range groups {
		groupsOut = append(groupsOut, dto.GroupFromModel(group))
	}

	usersOut := make([]*dto.User, 0, len(users))
	for _, user := range users {
		var role model.UserOrganizationRole
		orgAccess, ok := orgAccessesMap[user.ID]
		if ok {
			role = orgAccess.Role
		}
		usersOut = append(usersOut, dto.UserFromModel(user, nil, role))
	}

	userGroupsOut := make([]*dto.UserGroup, 0, len(userGroups))
	for _, userGroup := range userGroups {
		userGroupsOut = append(userGroupsOut, dto.UserGroupFromModel(userGroup))
	}

	return &dto.ListGroupsOutput{
		Groups:     groupsOut,
		Users:      usersOut,
		UserGroups: userGroupsOut,
	}, nil
}

func (s *serviceEE) Create(ctx context.Context, in dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
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

	group, err = s.Store.Group().Get(ctx, storeopts.GroupByID(group.ID))
	if err != nil {
		return nil, err
	}

	return &dto.CreateGroupOutput{
		Group: dto.GroupFromModel(group),
	}, nil
}

func (s *serviceEE) Update(ctx context.Context, in dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditGroup); err != nil {
		return nil, err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	group, err := s.Store.Group().Get(ctx, storeopts.GroupByID(groupID), storeopts.GroupByOrganizationID(currentOrg.ID))
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

		existingGroups, err := tx.User().ListGroups(ctx, storeopts.UserGroupByGroupID(group.ID))
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

	group, err = s.Store.Group().Get(ctx, storeopts.GroupByID(group.ID))
	if err != nil {
		return nil, err
	}

	return &dto.UpdateGroupOutput{
		Group: dto.GroupFromModel(group),
	}, nil
}

func (s *serviceEE) Delete(ctx context.Context, in dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error) {
	authorizer := authz.NewAuthorizer(s.Store)
	if err := authorizer.AuthorizeOperation(ctx, authz.OperationEditGroup); err != nil {
		return nil, err
	}

	currentOrg := ctxutils.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	group, err := s.Store.Group().Get(ctx, storeopts.GroupByID(groupID), storeopts.GroupByOrganizationID(currentOrg.ID))
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

	return &dto.DeleteGroupOutput{
		Group: dto.GroupFromModel(group),
	}, nil
}
