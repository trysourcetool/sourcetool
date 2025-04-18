package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type StoreOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isStoreOption()
}

func ByID(id uuid.UUID) StoreOption {
	return byIDOption{id: id}
}

type byIDOption struct {
	id uuid.UUID
}

func (o byIDOption) isStoreOption() {}

func (o byIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."id"`: o.id})
}

func ByEmail(email string) StoreOption {
	return byEmailOption{email: email}
}

type byEmailOption struct {
	email string
}

func (o byEmailOption) isStoreOption() {}

func (o byEmailOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."email"`: o.email})
}

func ByRefreshTokenHash(refreshTokenHash string) StoreOption {
	return byRefreshTokenHashOption{refreshTokenHash: refreshTokenHash}
}

type byRefreshTokenHashOption struct {
	refreshTokenHash string
}

func (o byRefreshTokenHashOption) isStoreOption() {}

func (o byRefreshTokenHashOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."refresh_token_hash"`: o.refreshTokenHash})
}

func ByOrganizationID(id uuid.UUID) StoreOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isStoreOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"user_organization_access" uoa ON u."id" = uoa."user_id"`).
		Where(sq.Eq{`uoa."organization_id"`: o.id})
}

func Limit(limit uint64) StoreOption {
	return limitOption{limit: limit}
}

type limitOption struct {
	limit uint64
}

func (o limitOption) isStoreOption() {}

func (o limitOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Limit(o.limit)
}

func Offset(offset uint64) StoreOption {
	return offsetOption{offset: offset}
}

type offsetOption struct {
	offset uint64
}

func (o offsetOption) isStoreOption() {}

func (o offsetOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Offset(o.offset)
}

func OrderBy(orderBy string) StoreOption {
	return orderByOption{orderBy: orderBy}
}

type orderByOption struct {
	orderBy string
}

func (o orderByOption) isStoreOption() {}

func (o orderByOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(o.orderBy)
}

type OrganizationAccessStoreOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isOrganizationAccessStoreOption()
}

func OrganizationAccessByUserID(id uuid.UUID) OrganizationAccessStoreOption {
	return organizationAccessByUserIDOption{id: id}
}

type organizationAccessByUserIDOption struct {
	id uuid.UUID
}

func (o organizationAccessByUserIDOption) isOrganizationAccessStoreOption() {}

func (o organizationAccessByUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."user_id"`: o.id})
}

func OrganizationAccessByUserIDs(ids []uuid.UUID) OrganizationAccessStoreOption {
	return organizationAccessByUserIDsOption{ids: ids}
}

type organizationAccessByUserIDsOption struct {
	ids []uuid.UUID
}

func (o organizationAccessByUserIDsOption) isOrganizationAccessStoreOption() {}

func (o organizationAccessByUserIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."user_id"`: o.ids})
}

func OrganizationAccessByOrganizationID(id uuid.UUID) OrganizationAccessStoreOption {
	return organizationAccessByOrganizationIDOption{id: id}
}

type organizationAccessByOrganizationIDOption struct {
	id uuid.UUID
}

func (o organizationAccessByOrganizationIDOption) isOrganizationAccessStoreOption() {}

func (o organizationAccessByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
		Where(sq.Eq{`o."id"`: o.id})
}

func OrganizationAccessByOrganizationSubdomain(subdomain string) OrganizationAccessStoreOption {
	return organizationAccessByOrganizationSubdomainOption{subdomain: subdomain}
}

type organizationAccessByOrganizationSubdomainOption struct {
	subdomain string
}

func (o organizationAccessByOrganizationSubdomainOption) isOrganizationAccessStoreOption() {}

func (o organizationAccessByOrganizationSubdomainOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
		Where(sq.Eq{`o."subdomain"`: o.subdomain})
}

func OrganizationAccessOrderBy(orderBy string) OrganizationAccessStoreOption {
	return organizationAccessOrderByOption{orderBy: orderBy}
}

type organizationAccessOrderByOption struct {
	orderBy string
}

func (o organizationAccessOrderByOption) isOrganizationAccessStoreOption() {}

func (o organizationAccessOrderByOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(o.orderBy)
}

func OrganizationAccessByRole(role UserOrganizationRole) OrganizationAccessStoreOption {
	return organizationAccessByRoleOption{role: role}
}

type organizationAccessByRoleOption struct {
	role UserOrganizationRole
}

func (o organizationAccessByRoleOption) isOrganizationAccessStoreOption() {}

func (o organizationAccessByRoleOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."role"`: o.role})
}

type GroupStoreOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isGroupStoreOption()
}

func GroupByUserID(id uuid.UUID) GroupStoreOption {
	return groupByUserIDOption{id: id}
}

type groupByUserIDOption struct {
	id uuid.UUID
}

func (o groupByUserIDOption) isGroupStoreOption() {}

func (o groupByUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ug."user_id"`: o.id})
}

func GroupByGroupID(id uuid.UUID) GroupStoreOption {
	return groupByGroupIDOption{id: id}
}

type groupByGroupIDOption struct {
	id uuid.UUID
}

func (o groupByGroupIDOption) isGroupStoreOption() {}

func (o groupByGroupIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ug."group_id"`: o.id})
}

func GroupByOrganizationID(id uuid.UUID) GroupStoreOption {
	return groupByOrganizationIDOption{id: id}
}

type groupByOrganizationIDOption struct {
	id uuid.UUID
}

func (o groupByOrganizationIDOption) isGroupStoreOption() {}

func (o groupByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"group" g ON g."id" = ug."group_id"`).
		Where(sq.Eq{`g."organization_id"`: o.id})
}

type InvitationStoreOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isInvitationStoreOption()
}

func InvitationByOrganizationID(id uuid.UUID) InvitationStoreOption {
	return invitationByOrganizationIDOption{id: id}
}

type invitationByOrganizationIDOption struct {
	id uuid.UUID
}

func (o invitationByOrganizationIDOption) isInvitationStoreOption() {}

func (o invitationByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."organization_id"`: o.id})
}

func InvitationByID(id uuid.UUID) InvitationStoreOption {
	return invitationByIDOption{id: id}
}

type invitationByIDOption struct {
	id uuid.UUID
}

func (o invitationByIDOption) isInvitationStoreOption() {}

func (o invitationByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."id"`: o.id})
}

func InvitationByEmail(email string) InvitationStoreOption {
	return invitationByEmailOption{email: email}
}

type invitationByEmailOption struct {
	email string
}

func (o invitationByEmailOption) isInvitationStoreOption() {}

func (o invitationByEmailOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."email"`: o.email})
}

type Store interface {
	Get(context.Context, ...StoreOption) (*User, error)
	List(context.Context, ...StoreOption) ([]*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	IsEmailExists(context.Context, string) (bool, error)

	GetOrganizationAccess(context.Context, ...OrganizationAccessStoreOption) (*UserOrganizationAccess, error)
	ListOrganizationAccesses(context.Context, ...OrganizationAccessStoreOption) ([]*UserOrganizationAccess, error)
	CreateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	UpdateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	DeleteOrganizationAccess(context.Context, *UserOrganizationAccess) error

	GetGroup(context.Context, ...GroupStoreOption) (*UserGroup, error)
	ListGroups(context.Context, ...GroupStoreOption) ([]*UserGroup, error)
	BulkInsertGroups(context.Context, []*UserGroup) error
	BulkDeleteGroups(context.Context, []*UserGroup) error

	GetInvitation(context.Context, ...InvitationStoreOption) (*UserInvitation, error)
	ListInvitations(context.Context, ...InvitationStoreOption) ([]*UserInvitation, error)
	DeleteInvitation(context.Context, *UserInvitation) error
	BulkInsertInvitations(context.Context, []*UserInvitation) error
	IsInvitationEmailExists(context.Context, uuid.UUID, string) (bool, error)
}
