package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/database"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

var _ database.UserStore = (*userStore)(nil)

type userStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newUserStore(db internal.DB) *userStore {
	return &userStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *userStore) Get(ctx context.Context, queries ...database.UserQuery) (*core.User, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.User{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *userStore) List(ctx context.Context, queries ...database.UserQuery) ([]*core.User, error) {
	query, args, err := s.buildQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.User, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *userStore) applyQueries(b sq.SelectBuilder, queries ...database.UserQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.UserByIDQuery:
			b = b.Where(sq.Eq{`u."id"`: q.ID})
		case database.UserByEmailQuery:
			b = b.Where(sq.Eq{`u."email"`: q.Email})
		case database.UserByRefreshTokenHashQuery:
			b = b.Where(sq.Eq{`u."refresh_token_hash"`: q.RefreshTokenHash})
		case database.UserByOrganizationIDQuery:
			b = b.
				InnerJoin(`"user_organization_access" uoa ON u."id" = uoa."user_id"`).
				Where(sq.Eq{`uoa."organization_id"`: q.OrganizationID})
		case database.UserLimitQuery:
			b = b.Limit(q.Limit)
		case database.UserOffsetQuery:
			b = b.Offset(q.Offset)
		case database.UserOrderByQuery:
			b = b.OrderBy(q.OrderBy)
		}
	}
	return b
}

func (s *userStore) buildQuery(ctx context.Context, queries ...database.UserQuery) (string, []any, error) {
	q := s.builder.Select(
		`u."id"`,
		`u."created_at"`,
		`u."email"`,
		`u."first_name"`,
		`u."last_name"`,
		`u."updated_at"`,
		`u."refresh_token_hash"`,
		`u."google_id"`,
	).
		From(`"user" u`)

	q = s.applyQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *userStore) Create(ctx context.Context, m *core.User) error {
	if _, err := s.builder.
		Insert(`"user"`).
		Columns(
			`"id"`,
			`"email"`,
			`"first_name"`,
			`"last_name"`,
			`"refresh_token_hash"`,
			`"google_id"`,
		).
		Values(
			m.ID,
			m.Email,
			m.FirstName,
			m.LastName,
			m.RefreshTokenHash,
			m.GoogleID,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) Update(ctx context.Context, m *core.User) error {
	if _, err := s.builder.
		Update(`"user"`).
		Set(`"email"`, m.Email).
		Set(`"first_name"`, m.FirstName).
		Set(`"last_name"`, m.LastName).
		Set(`"refresh_token_hash"`, m.RefreshTokenHash).
		Set(`"google_id"`, m.GoogleID).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) IsEmailExists(ctx context.Context, email string) (bool, error) {
	if _, err := s.Get(ctx, database.UserByEmail(email)); err != nil {
		if errdefs.IsUserNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *userStore) GetOrganizationAccess(ctx context.Context, queries ...database.UserOrganizationAccessQuery) (*core.UserOrganizationAccess, error) {
	query, args, err := s.buildOrganizationAccessQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.UserOrganizationAccess{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserOrganizationAccessNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *userStore) ListOrganizationAccesses(ctx context.Context, queries ...database.UserOrganizationAccessQuery) ([]*core.UserOrganizationAccess, error) {
	query, args, err := s.buildOrganizationAccessQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.UserOrganizationAccess, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *userStore) applyOrganizationAccessQueries(b sq.SelectBuilder, queries ...database.UserOrganizationAccessQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.UserOrganizationAccessByUserIDQuery:
			b = b.Where(sq.Eq{`uoa."user_id"`: q.UserID})
		case database.UserOrganizationAccessByUserIDsQuery:
			b = b.Where(sq.Eq{`uoa."user_id"`: q.UserIDs})
		case database.UserOrganizationAccessByOrganizationIDQuery:
			b = b.
				InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
				Where(sq.Eq{`o."id"`: q.OrganizationID})
		case database.UserOrganizationAccessByOrganizationSubdomainQuery:
			b = b.
				InnerJoin(`"organization" o ON o."id" = uoa."organization_id"`).
				Where(sq.Eq{`o."subdomain"`: q.OrganizationSubdomain})
		case database.UserOrganizationAccessByRoleQuery:
			b = b.Where(sq.Eq{`uoa."role"`: q.Role})
		}
	}
	return b
}

func (s *userStore) buildOrganizationAccessQuery(ctx context.Context, queries ...database.UserOrganizationAccessQuery) (string, []any, error) {
	q := s.builder.Select(
		`uoa."id"`,
		`uoa."user_id"`,
		`uoa."organization_id"`,
		`uoa."role"`,
		`uoa."created_at"`,
		`uoa."updated_at"`,
	).
		From(`"user_organization_access" uoa`)

	q = s.applyOrganizationAccessQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *userStore) CreateOrganizationAccess(ctx context.Context, m *core.UserOrganizationAccess) error {
	if _, err := s.builder.
		Insert(`"user_organization_access"`).
		Columns(
			`"id"`,
			`"user_id"`,
			`"organization_id"`,
			`"role"`,
		).
		Values(
			m.ID,
			m.UserID,
			m.OrganizationID,
			m.Role,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) UpdateOrganizationAccess(ctx context.Context, m *core.UserOrganizationAccess) error {
	if _, err := s.builder.
		Update(`"user_organization_access"`).
		Set(`"role"`, m.Role).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) DeleteOrganizationAccess(ctx context.Context, m *core.UserOrganizationAccess) error {
	if _, err := s.builder.
		Delete(`"user_organization_access"`).
		Where(sq.Eq{`"user_id"`: m.UserID, `"organization_id"`: m.OrganizationID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		if err == sql.ErrNoRows {
			return errdefs.ErrUserOrganizationAccessNotFound(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) GetGroup(ctx context.Context, queries ...database.UserGroupQuery) (*core.UserGroup, error) {
	query, args, err := s.buildGroupQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.UserGroup{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserGroupNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *userStore) ListGroups(ctx context.Context, queries ...database.UserGroupQuery) ([]*core.UserGroup, error) {
	query, args, err := s.buildGroupQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.UserGroup, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *userStore) applyGroupQueries(b sq.SelectBuilder, queries ...database.UserGroupQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.UserGroupByUserIDQuery:
			b = b.Where(sq.Eq{`ug."user_id"`: q.UserID})
		case database.UserGroupByGroupIDQuery:
			b = b.Where(sq.Eq{`ug."group_id"`: q.GroupID})
		case database.UserGroupByOrganizationIDQuery:
			b = b.
				InnerJoin(`"group" g ON g."id" = ug."group_id"`).
				Where(sq.Eq{`g."organization_id"`: q.OrganizationID})
		}
	}
	return b
}

func (s *userStore) buildGroupQuery(ctx context.Context, queries ...database.UserGroupQuery) (string, []any, error) {
	q := s.builder.Select(
		`ug."id"`,
		`ug."user_id"`,
		`ug."group_id"`,
		`ug."created_at"`,
		`ug."updated_at"`,
	).
		From(`"user_group" ug`)

	q = s.applyGroupQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *userStore) BulkInsertGroups(ctx context.Context, m []*core.UserGroup) error {
	if len(m) == 0 {
		return nil
	}

	q := s.builder.
		Insert(`"user_group"`).
		Columns(`"id"`, `"user_id"`, `"group_id"`)

	for _, v := range m {
		q = q.Values(v.ID, v.UserID, v.GroupID)
	}

	if _, err := q.
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) BulkDeleteGroups(ctx context.Context, m []*core.UserGroup) error {
	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *core.UserGroup, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := s.builder.
		Delete(`"user_group"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) GetInvitation(ctx context.Context, queries ...database.UserInvitationQuery) (*core.UserInvitation, error) {
	query, args, err := s.buildInvitationQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.UserInvitation{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserInvitationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *userStore) ListInvitations(ctx context.Context, queries ...database.UserInvitationQuery) ([]*core.UserInvitation, error) {
	query, args, err := s.buildInvitationQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.UserInvitation, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *userStore) applyInvitationQueries(b sq.SelectBuilder, queries ...database.UserInvitationQuery) sq.SelectBuilder {
	for _, q := range queries {
		switch q := q.(type) {
		case database.UserInvitationByOrganizationIDQuery:
			b = b.Where(sq.Eq{`ui."organization_id"`: q.OrganizationID})
		case database.UserInvitationByIDQuery:
			b = b.Where(sq.Eq{`ui."id"`: q.ID})
		case database.UserInvitationByEmailQuery:
			b = b.Where(sq.Eq{`ui."email"`: q.Email})
		}
	}
	return b
}

func (s *userStore) buildInvitationQuery(ctx context.Context, queries ...database.UserInvitationQuery) (string, []any, error) {
	q := s.builder.Select(
		`ui."id"`,
		`ui."organization_id"`,
		`ui."email"`,
		`ui."role"`,
		`ui."created_at"`,
		`ui."updated_at"`,
	).
		From(`"user_invitation" ui`)

	q = s.applyInvitationQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *userStore) DeleteInvitation(ctx context.Context, m *core.UserInvitation) error {
	if _, err := s.builder.
		Delete(`"user_invitation"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) BulkInsertInvitations(ctx context.Context, m []*core.UserInvitation) error {
	if len(m) == 0 {
		return nil
	}

	q := s.builder.
		Insert(`"user_invitation"`).
		Columns(`"id"`, `"organization_id"`, `"email"`, `"role"`)

	for _, v := range m {
		q = q.Values(v.ID, v.OrganizationID, v.Email, v.Role)
	}

	if _, err := q.
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) IsInvitationEmailExists(ctx context.Context, orgID uuid.UUID, email string) (bool, error) {
	if _, err := s.GetInvitation(ctx, database.UserInvitationByOrganizationID(orgID), database.UserInvitationByEmail(email)); err != nil {
		if errdefs.IsUserInvitationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
