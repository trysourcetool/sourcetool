package user

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Query interface{ isQuery() }

type ByIDQuery struct{ ID uuid.UUID }

func (q ByIDQuery) isQuery() {}

func ByID(id uuid.UUID) Query { return ByIDQuery{ID: id} }

type ByEmailQuery struct{ Email string }

func (q ByEmailQuery) isQuery() {}

func ByEmail(email string) Query { return ByEmailQuery{Email: email} }

type ByRefreshTokenHashQuery struct{ RefreshTokenHash string }

func (q ByRefreshTokenHashQuery) isQuery() {}

func ByRefreshTokenHash(refreshTokenHash string) Query {
	return ByRefreshTokenHashQuery{RefreshTokenHash: refreshTokenHash}
}

type ByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q ByOrganizationIDQuery) isQuery() {}

func ByOrganizationID(organizationID uuid.UUID) Query {
	return ByOrganizationIDQuery{OrganizationID: organizationID}
}

type LimitQuery struct{ Limit uint64 }

func (q LimitQuery) isQuery() {}

func Limit(limit uint64) Query { return LimitQuery{Limit: limit} }

type OffsetQuery struct{ Offset uint64 }

func (q OffsetQuery) isQuery() {}

func Offset(offset uint64) Query { return OffsetQuery{Offset: offset} }

type OrderByQuery struct{ OrderBy string }

func (q OrderByQuery) isQuery() {}

func OrderBy(orderBy string) Query { return OrderByQuery{OrderBy: orderBy} }

type OrganizationAccessQuery interface{ isOrganizationAccessQuery() }

type OrganizationAccessByUserIDQuery struct{ UserID uuid.UUID }

func (q OrganizationAccessByUserIDQuery) isOrganizationAccessQuery() {}

func OrganizationAccessByUserID(userID uuid.UUID) OrganizationAccessQuery {
	return OrganizationAccessByUserIDQuery{UserID: userID}
}

type OrganizationAccessByUserIDsQuery struct{ UserIDs []uuid.UUID }

func (q OrganizationAccessByUserIDsQuery) isOrganizationAccessQuery() {}

func OrganizationAccessByUserIDs(userIDs []uuid.UUID) OrganizationAccessQuery {
	return OrganizationAccessByUserIDsQuery{UserIDs: userIDs}
}

type OrganizationAccessByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q OrganizationAccessByOrganizationIDQuery) isOrganizationAccessQuery() {}

func OrganizationAccessByOrganizationID(organizationID uuid.UUID) OrganizationAccessQuery {
	return OrganizationAccessByOrganizationIDQuery{OrganizationID: organizationID}
}

type OrganizationAccessByOrganizationSubdomainQuery struct{ OrganizationSubdomain string }

func (q OrganizationAccessByOrganizationSubdomainQuery) isOrganizationAccessQuery() {}

func OrganizationAccessByOrganizationSubdomain(subdomain string) OrganizationAccessQuery {
	return OrganizationAccessByOrganizationSubdomainQuery{OrganizationSubdomain: subdomain}
}

type OrganizationAccessOrderByQuery struct{ OrderBy string }

func (q OrganizationAccessOrderByQuery) isOrganizationAccessQuery() {}

func OrganizationAccessOrderBy(orderBy string) OrganizationAccessQuery {
	return OrganizationAccessOrderByQuery{OrderBy: orderBy}
}

type OrganizationAccessByRoleQuery struct{ Role UserOrganizationRole }

func (q OrganizationAccessByRoleQuery) isOrganizationAccessQuery() {}

func OrganizationAccessByRole(role UserOrganizationRole) OrganizationAccessQuery {
	return OrganizationAccessByRoleQuery{Role: role}
}

type GroupQuery interface{ isGroupQuery() }

type GroupByUserIDQuery struct{ UserID uuid.UUID }

func (q GroupByUserIDQuery) isGroupQuery() {}

func GroupByUserID(userID uuid.UUID) GroupQuery { return GroupByUserIDQuery{UserID: userID} }

type GroupByGroupIDQuery struct{ GroupID uuid.UUID }

func (q GroupByGroupIDQuery) isGroupQuery() {}

func GroupByGroupID(groupID uuid.UUID) GroupQuery { return GroupByGroupIDQuery{GroupID: groupID} }

type GroupByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q GroupByOrganizationIDQuery) isGroupQuery() {}

func GroupByOrganizationID(organizationID uuid.UUID) GroupQuery {
	return GroupByOrganizationIDQuery{OrganizationID: organizationID}
}

type InvitationQuery interface{ isInvitationQuery() }

type InvitationByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q InvitationByOrganizationIDQuery) isInvitationQuery() {}

func InvitationByOrganizationID(organizationID uuid.UUID) InvitationQuery {
	return InvitationByOrganizationIDQuery{OrganizationID: organizationID}
}

type InvitationByIDQuery struct{ ID uuid.UUID }

func (q InvitationByIDQuery) isInvitationQuery() {}

func InvitationByID(id uuid.UUID) InvitationQuery { return InvitationByIDQuery{ID: id} }

type InvitationByEmailQuery struct{ Email string }

func (q InvitationByEmailQuery) isInvitationQuery() {}

func InvitationByEmail(email string) InvitationQuery { return InvitationByEmailQuery{Email: email} }

type Repository interface {
	Get(context.Context, ...Query) (*User, error)
	List(context.Context, ...Query) ([]*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	IsEmailExists(context.Context, string) (bool, error)

	GetOrganizationAccess(context.Context, ...OrganizationAccessQuery) (*UserOrganizationAccess, error)
	ListOrganizationAccesses(context.Context, ...OrganizationAccessQuery) ([]*UserOrganizationAccess, error)
	CreateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	UpdateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	DeleteOrganizationAccess(context.Context, *UserOrganizationAccess) error

	GetGroup(context.Context, ...GroupQuery) (*UserGroup, error)
	ListGroups(context.Context, ...GroupQuery) ([]*UserGroup, error)
	BulkInsertGroups(context.Context, []*UserGroup) error
	BulkDeleteGroups(context.Context, []*UserGroup) error

	GetInvitation(context.Context, ...InvitationQuery) (*UserInvitation, error)
	ListInvitations(context.Context, ...InvitationQuery) ([]*UserInvitation, error)
	DeleteInvitation(context.Context, *UserInvitation) error
	BulkInsertInvitations(context.Context, []*UserInvitation) error
	IsInvitationEmailExists(context.Context, uuid.UUID, string) (bool, error)
}
