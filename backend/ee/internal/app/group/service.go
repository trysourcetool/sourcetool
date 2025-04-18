package group

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/app/dto"
	groupSvc "github.com/trysourcetool/sourcetool/backend/internal/app/group"
	"github.com/trysourcetool/sourcetool/backend/internal/ctxutil"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/group"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
	"github.com/trysourcetool/sourcetool/backend/internal/permission"
	"github.com/trysourcetool/sourcetool/backend/pkg/errdefs"
	"github.com/trysourcetool/sourcetool/backend/pkg/ptrconv"
)

type serviceEE struct {
	*infra.Dependency
	*groupSvc.ServiceCE
}

func NewServiceEE(d *infra.Dependency) *serviceEE {
	return &serviceEE{
		Dependency: d,
		ServiceCE: groupSvc.NewServiceCE(
			infra.NewDependency(d.Repository, d.Mailer, d.PubSub),
		),
	}
}

func (s *serviceEE) Get(ctx context.Context, in dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	group, err := s.Repository.Group().Get(ctx, group.ByOrganizationID(currentOrg.ID), group.ByID(groupID))
	if err != nil {
		return nil, err
	}

	return &dto.GetGroupOutput{
		Group: dto.GroupFromModel(group),
	}, nil
}

func (s *serviceEE) List(ctx context.Context) (*dto.ListGroupsOutput, error) {
	currentOrg := ctxutil.CurrentOrganization(ctx)
	groups, err := s.Repository.Group().List(ctx, group.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	users, err := s.Repository.User().List(ctx, user.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userGroups, err := s.Repository.User().ListGroups(ctx, user.GroupByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	userIDs := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	orgAccesses, err := s.Repository.User().ListOrganizationAccesses(ctx, user.OrganizationAccessByUserIDs(userIDs))
	if err != nil {
		return nil, err
	}

	orgAccessesMap := make(map[uuid.UUID]*user.UserOrganizationAccess)
	for _, orgAccess := range orgAccesses {
		orgAccessesMap[orgAccess.UserID] = orgAccess
	}

	groupsOut := make([]*dto.Group, 0, len(groups))
	for _, group := range groups {
		groupsOut = append(groupsOut, dto.GroupFromModel(group))
	}

	usersOut := make([]*dto.User, 0, len(users))
	for _, u := range users {
		var role user.UserOrganizationRole
		orgAccess, ok := orgAccessesMap[u.ID]
		if ok {
			role = orgAccess.Role
		}
		usersOut = append(usersOut, dto.UserFromModel(u, nil, role))
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
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditGroup); err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)

	slugExists, err := s.Repository.Group().IsSlugExistsInOrganization(ctx, currentOrg.ID, in.Slug)
	if err != nil {
		return nil, err
	}
	if slugExists {
		return nil, errdefs.ErrGroupSlugAlreadyExists(errors.New("slug already exists"))
	}

	if !validateSlug(in.Slug) {
		return nil, errdefs.ErrInvalidArgument(errors.New("invalid slug"))
	}

	g := &group.Group{
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

	userGroups := make([]*user.UserGroup, 0, len(userIDs))
	for _, userID := range userIDs {
		userGroups = append(userGroups, &user.UserGroup{
			ID:      uuid.Must(uuid.NewV4()),
			UserID:  userID,
			GroupID: g.ID,
		})
	}

	if err := s.Repository.RunTransaction(func(tx db.Transaction) error {
		if err := tx.Group().Create(ctx, g); err != nil {
			return err
		}

		if err := tx.User().BulkInsertGroups(ctx, userGroups); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	g, err = s.Repository.Group().Get(ctx, group.ByID(g.ID))
	if err != nil {
		return nil, err
	}

	return &dto.CreateGroupOutput{
		Group: dto.GroupFromModel(g),
	}, nil
}

func (s *serviceEE) Update(ctx context.Context, in dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error) {
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditGroup); err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	g, err := s.Repository.Group().Get(ctx, group.ByID(groupID), group.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		g.Name = ptrconv.SafeValue(in.Name)
	}

	userIDs := make([]uuid.UUID, 0, len(in.UserIDs))
	for _, userID := range in.UserIDs {
		id, err := uuid.FromString(userID)
		if err != nil {
			return nil, errdefs.ErrInvalidArgument(err)
		}
		userIDs = append(userIDs, id)
	}

	userGroups := make([]*user.UserGroup, 0, len(userIDs))
	for _, userID := range userIDs {
		userGroups = append(userGroups, &user.UserGroup{
			ID:      uuid.Must(uuid.NewV4()),
			UserID:  userID,
			GroupID: g.ID,
		})
	}

	if err := s.Repository.RunTransaction(func(tx db.Transaction) error {
		if err := tx.Group().Update(ctx, g); err != nil {
			return err
		}

		existingGroups, err := tx.User().ListGroups(ctx, user.GroupByGroupID(g.ID))
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

	g, err = s.Repository.Group().Get(ctx, group.ByID(g.ID))
	if err != nil {
		return nil, err
	}

	return &dto.UpdateGroupOutput{
		Group: dto.GroupFromModel(g),
	}, nil
}

func (s *serviceEE) Delete(ctx context.Context, in dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error) {
	checker := permission.NewChecker(s.Repository)
	if err := checker.AuthorizeOperation(ctx, permission.OperationEditGroup); err != nil {
		return nil, err
	}

	currentOrg := ctxutil.CurrentOrganization(ctx)
	groupID, err := uuid.FromString(in.GroupID)
	if err != nil {
		return nil, errdefs.ErrInvalidArgument(err)
	}

	g, err := s.Repository.Group().Get(ctx, group.ByID(groupID), group.ByOrganizationID(currentOrg.ID))
	if err != nil {
		return nil, err
	}

	if err := s.Repository.RunTransaction(func(tx db.Transaction) error {
		if err := tx.Group().Delete(ctx, g); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.DeleteGroupOutput{
		Group: dto.GroupFromModel(g),
	}, nil
}
