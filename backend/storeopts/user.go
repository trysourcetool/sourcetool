package storeopts

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
)

type UserOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func UserByID(id uuid.UUID) UserOption {
	return userByIDOption{id: id}
}

type userByIDOption struct {
	id uuid.UUID
}

func (o userByIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."id"`: o.id})
}

func UserByEmail(email string) UserOption {
	return userByEmailOption{email: email}
}

type userByEmailOption struct {
	email string
}

func (o userByEmailOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."email"`: o.email})
}

func UserBySecret(secret string) UserOption {
	return userBySecretOption{secret: secret}
}

type userBySecretOption struct {
	secret string
}

func (o userBySecretOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."secret"`: o.secret})
}

func UserByOrganizationID(id uuid.UUID) UserOption {
	return userByOrganizationIDOption{id: id}
}

type userByOrganizationIDOption struct {
	id uuid.UUID
}

func (o userByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"user_organization_access" uoa ON u."id" = uoa."user_id"`).
		Where(sq.Eq{`uoa."organization_id"`: o.id})
}

func UserLimit(limit uint64) UserOption {
	return userLimitOption{limit: limit}
}

type userLimitOption struct {
	limit uint64
}

func (o userLimitOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Limit(o.limit)
}

func UserOffset(offset uint64) UserOption {
	return userOffsetOption{offset: offset}
}

type userOffsetOption struct {
	offset uint64
}

func (o userOffsetOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Offset(o.offset)
}

func UserOrderBy(orderBy string) UserOption {
	return userOrderByOption{orderBy: orderBy}
}

type userOrderByOption struct {
	orderBy string
}

func (o userOrderByOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(o.orderBy)
}

type UserRegistrationRequestOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func UserRegistrationRequestByEmail(email string) UserRegistrationRequestOption {
	return userRegistrationRequestByEmailOption{email: email}
}

type userRegistrationRequestByEmailOption struct {
	email string
}

func (o userRegistrationRequestByEmailOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`urr."email"`: o.email})
}

type UserOrganizationAccessOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func UserOrganizationAccessByUserID(id uuid.UUID) UserOrganizationAccessOption {
	return userOrganizationAccessByUserIDOption{id: id}
}

type userOrganizationAccessByUserIDOption struct {
	id uuid.UUID
}

func (o userOrganizationAccessByUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."user_id"`: o.id})
}

func UserOrganizationAccessByUserIDs(ids []uuid.UUID) UserOrganizationAccessOption {
	return userOrganizationAccessByUserIDsOption{ids: ids}
}

type userOrganizationAccessByUserIDsOption struct {
	ids []uuid.UUID
}

func (o userOrganizationAccessByUserIDsOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."user_id"`: o.ids})
}

func UserOrganizationAccessByOrganizationID(id uuid.UUID) UserOrganizationAccessOption {
	return userOrganizationAccessByOrganizationIDOption{id: id}
}

type userOrganizationAccessByOrganizationIDOption struct {
	id uuid.UUID
}

func (o userOrganizationAccessByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
		Where(sq.Eq{`o."id"`: o.id})
}

func UserOrganizationAccessByOrganizationSubdomain(subdomain string) UserOrganizationAccessOption {
	return userOrganizationAccessByOrganizationSubdomainOption{subdomain: subdomain}
}

type userOrganizationAccessByOrganizationSubdomainOption struct {
	subdomain string
}

func (o userOrganizationAccessByOrganizationSubdomainOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
		Where(sq.Eq{`o."subdomain"`: o.subdomain})
}

type UserGroupOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func UserGroupByUserID(id uuid.UUID) UserGroupOption {
	return userGroupByUserIDOption{id: id}
}

type userGroupByUserIDOption struct {
	id uuid.UUID
}

func (o userGroupByUserIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ug."user_id"`: o.id})
}

func UserGroupByGroupID(id uuid.UUID) UserGroupOption {
	return userGroupByGroupIDOption{id: id}
}

type userGroupByGroupIDOption struct {
	id uuid.UUID
}

func (o userGroupByGroupIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ug."group_id"`: o.id})
}

func UserGroupByOrganizationID(id uuid.UUID) UserGroupOption {
	return userGroupByOrganizationIDOption{id: id}
}

type userGroupByOrganizationIDOption struct {
	id uuid.UUID
}

func (o userGroupByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"group" g ON g."id" = ug."group_id"`).
		Where(sq.Eq{`g."organization_id"`: o.id})
}

type UserInvitationOption interface {
	Apply(sq.SelectBuilder) sq.SelectBuilder
}

func UserInvitationByOrganizationID(id uuid.UUID) UserInvitationOption {
	return userInvitationByOrganizationIDOption{id: id}
}

type userInvitationByOrganizationIDOption struct {
	id uuid.UUID
}

func (o userInvitationByOrganizationIDOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."organization_id"`: o.id})
}

func UserInvitationByEmail(email string) UserInvitationOption {
	return userInvitationByEmailOption{email: email}
}

type userInvitationByEmailOption struct {
	email string
}

func (o userInvitationByEmailOption) Apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."email"`: o.email})
}
