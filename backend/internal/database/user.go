package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type UserQuery interface{ isUserQuery() }

type UserByIDQuery struct{ ID uuid.UUID }

func (q UserByIDQuery) isUserQuery() {}

func UserByID(id uuid.UUID) UserQuery { return UserByIDQuery{ID: id} }

type UserByEmailQuery struct{ Email string }

func (q UserByEmailQuery) isUserQuery() {}

func UserByEmail(email string) UserQuery { return UserByEmailQuery{Email: email} }

type UserByRefreshTokenHashQuery struct{ RefreshTokenHash string }

func (q UserByRefreshTokenHashQuery) isUserQuery() {}

func UserByRefreshTokenHash(refreshTokenHash string) UserQuery {
	return UserByRefreshTokenHashQuery{RefreshTokenHash: refreshTokenHash}
}

type UserByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q UserByOrganizationIDQuery) isUserQuery() {}

func UserByOrganizationID(organizationID uuid.UUID) UserQuery {
	return UserByOrganizationIDQuery{OrganizationID: organizationID}
}

type UserLimitQuery struct{ Limit uint64 }

func (q UserLimitQuery) isUserQuery() {}

func UserLimit(limit uint64) UserQuery { return UserLimitQuery{Limit: limit} }

type UserOffsetQuery struct{ Offset uint64 }

func (q UserOffsetQuery) isUserQuery() {}

func UserOffset(offset uint64) UserQuery { return UserOffsetQuery{Offset: offset} }

type UserOrderByQuery struct{ OrderBy string }

func (q UserOrderByQuery) isUserQuery() {}

func UserOrderBy(orderBy string) UserQuery { return UserOrderByQuery{OrderBy: orderBy} }

type UserOrganizationAccessQuery interface{ isUserOrganizationAccessQuery() }

type UserOrganizationAccessByUserIDQuery struct{ UserID uuid.UUID }

func (q UserOrganizationAccessByUserIDQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByUserID(userID uuid.UUID) UserOrganizationAccessQuery {
	return UserOrganizationAccessByUserIDQuery{UserID: userID}
}

type UserOrganizationAccessByUserIDsQuery struct{ UserIDs []uuid.UUID }

func (q UserOrganizationAccessByUserIDsQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByUserIDs(userIDs []uuid.UUID) UserOrganizationAccessQuery {
	return UserOrganizationAccessByUserIDsQuery{UserIDs: userIDs}
}

type UserOrganizationAccessByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q UserOrganizationAccessByOrganizationIDQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByOrganizationID(organizationID uuid.UUID) UserOrganizationAccessQuery {
	return UserOrganizationAccessByOrganizationIDQuery{OrganizationID: organizationID}
}

type UserOrganizationAccessByOrganizationSubdomainQuery struct{ OrganizationSubdomain string }

func (q UserOrganizationAccessByOrganizationSubdomainQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByOrganizationSubdomain(subdomain string) UserOrganizationAccessQuery {
	return UserOrganizationAccessByOrganizationSubdomainQuery{OrganizationSubdomain: subdomain}
}

type UserOrganizationAccessOrderByQuery struct{ OrderBy string }

func (q UserOrganizationAccessOrderByQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessOrderBy(orderBy string) UserOrganizationAccessQuery {
	return UserOrganizationAccessOrderByQuery{OrderBy: orderBy}
}

type UserOrganizationAccessByRoleQuery struct{ Role core.UserOrganizationRole }

func (q UserOrganizationAccessByRoleQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByRole(role core.UserOrganizationRole) UserOrganizationAccessQuery {
	return UserOrganizationAccessByRoleQuery{Role: role}
}

type UserGroupQuery interface{ isUserGroupQuery() }

type UserGroupByUserIDQuery struct{ UserID uuid.UUID }

func (q UserGroupByUserIDQuery) isUserGroupQuery() {}

func UserGroupByUserID(userID uuid.UUID) UserGroupQuery {
	return UserGroupByUserIDQuery{UserID: userID}
}

type UserGroupByGroupIDQuery struct{ GroupID uuid.UUID }

func (q UserGroupByGroupIDQuery) isUserGroupQuery() {}

func UserGroupByGroupID(groupID uuid.UUID) UserGroupQuery {
	return UserGroupByGroupIDQuery{GroupID: groupID}
}

type UserGroupByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q UserGroupByOrganizationIDQuery) isUserGroupQuery() {}

func UserGroupByOrganizationID(organizationID uuid.UUID) UserGroupQuery {
	return UserGroupByOrganizationIDQuery{OrganizationID: organizationID}
}

type UserInvitationQuery interface{ isUserInvitationQuery() }

type UserInvitationByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q UserInvitationByOrganizationIDQuery) isUserInvitationQuery() {}

func UserInvitationByOrganizationID(organizationID uuid.UUID) UserInvitationQuery {
	return UserInvitationByOrganizationIDQuery{OrganizationID: organizationID}
}

type UserInvitationByIDQuery struct{ ID uuid.UUID }

func (q UserInvitationByIDQuery) isUserInvitationQuery() {}

func UserInvitationByID(id uuid.UUID) UserInvitationQuery {
	return UserInvitationByIDQuery{ID: id}
}

type UserInvitationByEmailQuery struct{ Email string }

func (q UserInvitationByEmailQuery) isUserInvitationQuery() {}

func UserInvitationByEmail(email string) UserInvitationQuery {
	return UserInvitationByEmailQuery{Email: email}
}

type UserStore interface {
	Get(ctx context.Context, queries ...UserQuery) (*core.User, error)
	List(ctx context.Context, queries ...UserQuery) ([]*core.User, error)
	Create(ctx context.Context, m *core.User) error
	Update(ctx context.Context, m *core.User) error
	IsEmailExists(ctx context.Context, email string) (bool, error)

	GetOrganizationAccess(ctx context.Context, queries ...UserOrganizationAccessQuery) (*core.UserOrganizationAccess, error)
	ListOrganizationAccesses(ctx context.Context, queries ...UserOrganizationAccessQuery) ([]*core.UserOrganizationAccess, error)
	CreateOrganizationAccess(ctx context.Context, m *core.UserOrganizationAccess) error
	UpdateOrganizationAccess(ctx context.Context, m *core.UserOrganizationAccess) error
	DeleteOrganizationAccess(ctx context.Context, m *core.UserOrganizationAccess) error

	GetGroup(ctx context.Context, queries ...UserGroupQuery) (*core.UserGroup, error)
	ListGroups(ctx context.Context, queries ...UserGroupQuery) ([]*core.UserGroup, error)
	BulkInsertGroups(ctx context.Context, m []*core.UserGroup) error
	BulkDeleteGroups(ctx context.Context, m []*core.UserGroup) error

	GetInvitation(ctx context.Context, queries ...UserInvitationQuery) (*core.UserInvitation, error)
	ListInvitations(ctx context.Context, queries ...UserInvitationQuery) ([]*core.UserInvitation, error)
	DeleteInvitation(ctx context.Context, m *core.UserInvitation) error
	BulkInsertInvitations(ctx context.Context, m []*core.UserInvitation) error
	IsInvitationEmailExists(ctx context.Context, orgID uuid.UUID, email string) (bool, error)
}
