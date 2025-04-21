package postgres

import (
	sq "github.com/Masterminds/squirrel"
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

func applyUserQueries(b sq.SelectBuilder, queries ...UserQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case UserByIDQuery:
			b = b.Where(sq.Eq{`u."id"`: q.ID})
		case UserByEmailQuery:
			b = b.Where(sq.Eq{`u."email"`: q.Email})
		case UserByRefreshTokenHashQuery:
			b = b.Where(sq.Eq{`u."refresh_token_hash"`: q.RefreshTokenHash})
		case UserByOrganizationIDQuery:
			b = b.
				InnerJoin(`"user_organization_access" uoa ON u."id" = uoa."user_id"`).
				Where(sq.Eq{`uoa."organization_id"`: q.OrganizationID})
		case UserLimitQuery:
			b = b.Limit(q.Limit)
		case UserOffsetQuery:
			b = b.Offset(q.Offset)
		case UserOrderByQuery:
			b = b.OrderBy(q.OrderBy)
		}
	}
	return b
}

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

func applyUserOrganizationAccessQueries(b sq.SelectBuilder, queries ...UserOrganizationAccessQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case UserOrganizationAccessByUserIDQuery:
			b = b.Where(sq.Eq{`uoa."user_id"`: q.UserID})
		case UserOrganizationAccessByUserIDsQuery:
			b = b.Where(sq.Eq{`uoa."user_id"`: q.UserIDs})
		case UserOrganizationAccessByOrganizationIDQuery:
			b = b.
				InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
				Where(sq.Eq{`o."id"`: q.OrganizationID})
		case UserOrganizationAccessByOrganizationSubdomainQuery:
			b = b.
				InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
				Where(sq.Eq{`o."subdomain"`: q.OrganizationSubdomain})
		case UserOrganizationAccessByRoleQuery:
			b = b.Where(sq.Eq{`uoa."role"`: q.Role})
		}
	}
	return b
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

func applyUserGroupQueries(b sq.SelectBuilder, queries ...UserGroupQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case UserGroupByUserIDQuery:
			b = b.Where(sq.Eq{`ug."user_id"`: q.UserID})
		case UserGroupByGroupIDQuery:
			b = b.Where(sq.Eq{`ug."group_id"`: q.GroupID})
		case UserGroupByOrganizationIDQuery:
			b = b.
				InnerJoin(`"group" g ON g."id" = ug."group_id"`).
				Where(sq.Eq{`g."organization_id"`: q.OrganizationID})
		}
	}
	return b
}

type UserInvitationQuery interface{ isUserInvitationQuery() }

type UserInvitationByOrganizationIDQuery struct{ OrganizationID uuid.UUID }

func (q UserInvitationByOrganizationIDQuery) isUserInvitationQuery() {}

func UserInvitationByOrganizationID(organizationID uuid.UUID) UserInvitationQuery {
	return UserInvitationByOrganizationIDQuery{OrganizationID: organizationID}
}

type UserInvitationByIDQuery struct{ ID uuid.UUID }

func (q UserInvitationByIDQuery) isUserInvitationQuery() {}

func UserInvitationByID(id uuid.UUID) UserInvitationQuery { return UserInvitationByIDQuery{ID: id} }

type UserInvitationByEmailQuery struct{ Email string }

func (q UserInvitationByEmailQuery) isUserInvitationQuery() {}

func UserInvitationByEmail(email string) UserInvitationQuery {
	return UserInvitationByEmailQuery{Email: email}
}

func applyUserInvitationQueries(b sq.SelectBuilder, queries ...UserInvitationQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case UserInvitationByOrganizationIDQuery:
			b = b.Where(sq.Eq{`ui."organization_id"`: q.OrganizationID})
		case UserInvitationByIDQuery:
			b = b.Where(sq.Eq{`ui."id"`: q.ID})
		case UserInvitationByEmailQuery:
			b = b.Where(sq.Eq{`ui."email"`: q.Email})
		}
	}
	return b
}
