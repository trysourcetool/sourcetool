package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
)

type UserQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isUserQuery()
}

type userByIDQuery struct{ id uuid.UUID }

func (q userByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."id"`: q.id})
}

func (userByIDQuery) isUserQuery() {}

func UserByID(id uuid.UUID) UserQuery { return userByIDQuery{id: id} }

type userByEmailQuery struct{ email string }

func (q userByEmailQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."email"`: q.email})
}

func (userByEmailQuery) isUserQuery() {}

func UserByEmail(email string) UserQuery { return userByEmailQuery{email: email} }

type userByRefreshTokenHashQuery struct{ refreshTokenHash string }

func (q userByRefreshTokenHashQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`u."refresh_token_hash"`: q.refreshTokenHash})
}

func (userByRefreshTokenHashQuery) isUserQuery() {}

func UserByRefreshTokenHash(refreshTokenHash string) UserQuery {
	return userByRefreshTokenHashQuery{refreshTokenHash: refreshTokenHash}
}

type userByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q userByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"user_organization_access" uoa ON u."id" = uoa."user_id"`).
		Where(sq.Eq{`uoa."organization_id"`: q.organizationID})
}

func (userByOrganizationIDQuery) isUserQuery() {}

func UserByOrganizationID(organizationID uuid.UUID) UserQuery {
	return userByOrganizationIDQuery{organizationID: organizationID}
}

type userLimitQuery struct{ limit uint64 }

func (q userLimitQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Limit(q.limit)
}

func (userLimitQuery) isUserQuery() {}

func UserLimit(limit uint64) UserQuery { return userLimitQuery{limit: limit} }

type userOffsetQuery struct{ offset uint64 }

func (q userOffsetQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Offset(q.offset)
}

func (userOffsetQuery) isUserQuery() {}

func UserOffset(offset uint64) UserQuery { return userOffsetQuery{offset: offset} }

type userOrderByQuery struct{ orderBy string }

func (q userOrderByQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(q.orderBy)
}

func (userOrderByQuery) isUserQuery() {}

func UserOrderBy(orderBy string) UserQuery { return userOrderByQuery{orderBy: orderBy} }

type UserOrganizationAccessQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isUserOrganizationAccessQuery()
}

type userOrganizationAccessByUserIDQuery struct{ userID uuid.UUID }

func (q userOrganizationAccessByUserIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."user_id"`: q.userID})
}

func (userOrganizationAccessByUserIDQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByUserID(userID uuid.UUID) UserOrganizationAccessQuery {
	return userOrganizationAccessByUserIDQuery{userID: userID}
}

type userOrganizationAccessByUserIDsQuery struct{ userIDs []uuid.UUID }

func (q userOrganizationAccessByUserIDsQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."user_id"`: q.userIDs})
}

func (userOrganizationAccessByUserIDsQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByUserIDs(userIDs []uuid.UUID) UserOrganizationAccessQuery {
	return userOrganizationAccessByUserIDsQuery{userIDs: userIDs}
}

type userOrganizationAccessByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q userOrganizationAccessByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."organization_id"`: q.organizationID})
}

func (userOrganizationAccessByOrganizationIDQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByOrganizationID(organizationID uuid.UUID) UserOrganizationAccessQuery {
	return userOrganizationAccessByOrganizationIDQuery{organizationID: organizationID}
}

type userOrganizationAccessByOrganizationSubdomainQuery struct{ organizationSubdomain string }

func (q userOrganizationAccessByOrganizationSubdomainQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
		Where(sq.Eq{`o."subdomain"`: q.organizationSubdomain})
}

func (userOrganizationAccessByOrganizationSubdomainQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByOrganizationSubdomain(subdomain string) UserOrganizationAccessQuery {
	return userOrganizationAccessByOrganizationSubdomainQuery{organizationSubdomain: subdomain}
}

type userOrganizationAccessOrderByQuery struct{ orderBy string }

func (q userOrganizationAccessOrderByQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(q.orderBy)
}

func (userOrganizationAccessOrderByQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessOrderBy(orderBy string) UserOrganizationAccessQuery {
	return userOrganizationAccessOrderByQuery{orderBy: orderBy}
}

type userOrganizationAccessByRoleQuery struct{ role core.UserOrganizationRole }

func (q userOrganizationAccessByRoleQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`uoa."role"`: q.role})
}

func (userOrganizationAccessByRoleQuery) isUserOrganizationAccessQuery() {}

func UserOrganizationAccessByRole(role core.UserOrganizationRole) UserOrganizationAccessQuery {
	return userOrganizationAccessByRoleQuery{role: role}
}

type UserGroupQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isUserGroupQuery()
}

type userGroupByUserIDQuery struct{ userID uuid.UUID }

func (q userGroupByUserIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ug."user_id"`: q.userID})
}

func (userGroupByUserIDQuery) isUserGroupQuery() {}

func UserGroupByUserID(userID uuid.UUID) UserGroupQuery {
	return userGroupByUserIDQuery{userID: userID}
}

type userGroupByGroupIDQuery struct{ groupID uuid.UUID }

func (q userGroupByGroupIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ug."group_id"`: q.groupID})
}

func (userGroupByGroupIDQuery) isUserGroupQuery() {}

func UserGroupByGroupID(groupID uuid.UUID) UserGroupQuery {
	return userGroupByGroupIDQuery{groupID: groupID}
}

type userGroupByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q userGroupByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.
		InnerJoin(`"group" g ON g."id" = ug."group_id"`).
		Where(sq.Eq{`g."organization_id"`: q.organizationID})
}

func (userGroupByOrganizationIDQuery) isUserGroupQuery() {}

func UserGroupByOrganizationID(organizationID uuid.UUID) UserGroupQuery {
	return userGroupByOrganizationIDQuery{organizationID: organizationID}
}

type UserInvitationQuery interface {
	apply(b sq.SelectBuilder) sq.SelectBuilder
	isUserInvitationQuery()
}

type userInvitationByOrganizationIDQuery struct{ organizationID uuid.UUID }

func (q userInvitationByOrganizationIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."organization_id"`: q.organizationID})
}

func (userInvitationByOrganizationIDQuery) isUserInvitationQuery() {}

func UserInvitationByOrganizationID(organizationID uuid.UUID) UserInvitationQuery {
	return userInvitationByOrganizationIDQuery{organizationID: organizationID}
}

type userInvitationByIDQuery struct{ id uuid.UUID }

func (q userInvitationByIDQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."id"`: q.id})
}

func (userInvitationByIDQuery) isUserInvitationQuery() {}

func UserInvitationByID(id uuid.UUID) UserInvitationQuery { return userInvitationByIDQuery{id: id} }

type userInvitationByEmailQuery struct{ email string }

func (q userInvitationByEmailQuery) apply(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(sq.Eq{`ui."email"`: q.email})
}

func (userInvitationByEmailQuery) isUserInvitationQuery() {}

func UserInvitationByEmail(email string) UserInvitationQuery {
	return userInvitationByEmailQuery{email: email}
}
