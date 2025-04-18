package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type RepositoryOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isRepositoryOption()
}

func ByID(id uuid.UUID) RepositoryOption {
	return byIDOption{id: id}
}

type byIDOption struct {
	id uuid.UUID
}

func (o byIDOption) isRepositoryOption() {}

func (o byIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."id"`: o.id})
}

func ByEmail(email string) RepositoryOption {
	return byEmailOption{email: email}
}

type byEmailOption struct {
	email string
}

func (o byEmailOption) isRepositoryOption() {}

func (o byEmailOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."email"`: o.email})
}

func ByRefreshTokenHash(refreshTokenHash string) RepositoryOption {
	return byRefreshTokenHashOption{refreshTokenHash: refreshTokenHash}
}

type byRefreshTokenHashOption struct {
	refreshTokenHash string
}

func (o byRefreshTokenHashOption) isRepositoryOption() {}

func (o byRefreshTokenHashOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."refresh_token_hash"`: o.refreshTokenHash})
}

func ByOrganizationID(id uuid.UUID) RepositoryOption {
	return byOrganizationIDOption{id: id}
}

type byOrganizationIDOption struct {
	id uuid.UUID
}

func (o byOrganizationIDOption) isRepositoryOption() {}

func (o byOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"user_organization_access" uoa ON u."id" = uoa."user_id"`).
		Where(sq.Eq{`uoa."organization_id"`: o.id})
}

func Limit(limit uint64) RepositoryOption {
	return limitOption{limit: limit}
}

type limitOption struct {
	limit uint64
}

func (o limitOption) isRepositoryOption() {}

func (o limitOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Limit(o.limit)
}

func Offset(offset uint64) RepositoryOption {
	return offsetOption{offset: offset}
}

type offsetOption struct {
	offset uint64
}

func (o offsetOption) isRepositoryOption() {}

func (o offsetOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Offset(o.offset)
}

func OrderBy(orderBy string) RepositoryOption {
	return orderByOption{orderBy: orderBy}
}

type orderByOption struct {
	orderBy string
}

func (o orderByOption) isRepositoryOption() {}

func (o orderByOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(o.orderBy)
}

type OrganizationAccessRepositoryOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isOrganizationAccessRepositoryOption()
}

func OrganizationAccessByUserID(id uuid.UUID) OrganizationAccessRepositoryOption {
	return organizationAccessByUserIDOption{id: id}
}

type organizationAccessByUserIDOption struct {
	id uuid.UUID
}

func (o organizationAccessByUserIDOption) isOrganizationAccessRepositoryOption() {}

func (o organizationAccessByUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."user_id"`: o.id})
}

func OrganizationAccessByUserIDs(ids []uuid.UUID) OrganizationAccessRepositoryOption {
	return organizationAccessByUserIDsOption{ids: ids}
}

type organizationAccessByUserIDsOption struct {
	ids []uuid.UUID
}

func (o organizationAccessByUserIDsOption) isOrganizationAccessRepositoryOption() {}

func (o organizationAccessByUserIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."user_id"`: o.ids})
}

func OrganizationAccessByOrganizationID(id uuid.UUID) OrganizationAccessRepositoryOption {
	return organizationAccessByOrganizationIDOption{id: id}
}

type organizationAccessByOrganizationIDOption struct {
	id uuid.UUID
}

func (o organizationAccessByOrganizationIDOption) isOrganizationAccessRepositoryOption() {}

func (o organizationAccessByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
		Where(sq.Eq{`o."id"`: o.id})
}

func OrganizationAccessByOrganizationSubdomain(subdomain string) OrganizationAccessRepositoryOption {
	return organizationAccessByOrganizationSubdomainOption{subdomain: subdomain}
}

type organizationAccessByOrganizationSubdomainOption struct {
	subdomain string
}

func (o organizationAccessByOrganizationSubdomainOption) isOrganizationAccessRepositoryOption() {}

func (o organizationAccessByOrganizationSubdomainOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
		Where(sq.Eq{`o."subdomain"`: o.subdomain})
}

func OrganizationAccessOrderBy(orderBy string) OrganizationAccessRepositoryOption {
	return organizationAccessOrderByOption{orderBy: orderBy}
}

type organizationAccessOrderByOption struct {
	orderBy string
}

func (o organizationAccessOrderByOption) isOrganizationAccessRepositoryOption() {}

func (o organizationAccessOrderByOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(o.orderBy)
}

func OrganizationAccessByRole(role UserOrganizationRole) OrganizationAccessRepositoryOption {
	return organizationAccessByRoleOption{role: role}
}

type organizationAccessByRoleOption struct {
	role UserOrganizationRole
}

func (o organizationAccessByRoleOption) isOrganizationAccessRepositoryOption() {}

func (o organizationAccessByRoleOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."role"`: o.role})
}

type GroupRepositoryOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isGroupRepositoryOption()
}

func GroupByUserID(id uuid.UUID) GroupRepositoryOption {
	return groupByUserIDOption{id: id}
}

type groupByUserIDOption struct {
	id uuid.UUID
}

func (o groupByUserIDOption) isGroupRepositoryOption() {}

func (o groupByUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ug."user_id"`: o.id})
}

func GroupByGroupID(id uuid.UUID) GroupRepositoryOption {
	return groupByGroupIDOption{id: id}
}

type groupByGroupIDOption struct {
	id uuid.UUID
}

func (o groupByGroupIDOption) isGroupRepositoryOption() {}

func (o groupByGroupIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ug."group_id"`: o.id})
}

func GroupByOrganizationID(id uuid.UUID) GroupRepositoryOption {
	return groupByOrganizationIDOption{id: id}
}

type groupByOrganizationIDOption struct {
	id uuid.UUID
}

func (o groupByOrganizationIDOption) isGroupRepositoryOption() {}

func (o groupByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"group" g ON g."id" = ug."group_id"`).
		Where(sq.Eq{`g."organization_id"`: o.id})
}

type InvitationRepositoryOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
	isInvitationRepositoryOption()
}

func InvitationByOrganizationID(id uuid.UUID) InvitationRepositoryOption {
	return invitationByOrganizationIDOption{id: id}
}

type invitationByOrganizationIDOption struct {
	id uuid.UUID
}

func (o invitationByOrganizationIDOption) isInvitationRepositoryOption() {}

func (o invitationByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."organization_id"`: o.id})
}

func InvitationByID(id uuid.UUID) InvitationRepositoryOption {
	return invitationByIDOption{id: id}
}

type invitationByIDOption struct {
	id uuid.UUID
}

func (o invitationByIDOption) isInvitationRepositoryOption() {}

func (o invitationByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."id"`: o.id})
}

func InvitationByEmail(email string) InvitationRepositoryOption {
	return invitationByEmailOption{email: email}
}

type invitationByEmailOption struct {
	email string
}

func (o invitationByEmailOption) isInvitationRepositoryOption() {}

func (o invitationByEmailOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."email"`: o.email})
}

type Repository interface {
	Get(context.Context, ...RepositoryOption) (*User, error)
	List(context.Context, ...RepositoryOption) ([]*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	IsEmailExists(context.Context, string) (bool, error)

	GetOrganizationAccess(context.Context, ...OrganizationAccessRepositoryOption) (*UserOrganizationAccess, error)
	ListOrganizationAccesses(context.Context, ...OrganizationAccessRepositoryOption) ([]*UserOrganizationAccess, error)
	CreateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	UpdateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	DeleteOrganizationAccess(context.Context, *UserOrganizationAccess) error

	GetGroup(context.Context, ...GroupRepositoryOption) (*UserGroup, error)
	ListGroups(context.Context, ...GroupRepositoryOption) ([]*UserGroup, error)
	BulkInsertGroups(context.Context, []*UserGroup) error
	BulkDeleteGroups(context.Context, []*UserGroup) error

	GetInvitation(context.Context, ...InvitationRepositoryOption) (*UserInvitation, error)
	ListInvitations(context.Context, ...InvitationRepositoryOption) ([]*UserInvitation, error)
	DeleteInvitation(context.Context, *UserInvitation) error
	BulkInsertInvitations(context.Context, []*UserInvitation) error
	IsInvitationEmailExists(context.Context, uuid.UUID, string) (bool, error)
}
