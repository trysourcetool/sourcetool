package user

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/internal/domain/user"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db"
)

type RepositoryCE struct {
	db      db.DB
	builder sq.StatementBuilderType
}

func NewRepositoryCE(db db.DB) *RepositoryCE {
	return &RepositoryCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RepositoryCE) Get(ctx context.Context, opts ...user.RepositoryOption) (*user.User, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := user.User{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) List(ctx context.Context, opts ...user.RepositoryOption) ([]*user.User, error) {
	query, args, err := r.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*user.User, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *RepositoryCE) buildQuery(ctx context.Context, opts ...user.RepositoryOption) (string, []any, error) {
	q := r.builder.Select(
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

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (r *RepositoryCE) Create(ctx context.Context, m *user.User) error {
	if _, err := r.builder.
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
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) Update(ctx context.Context, m *user.User) error {
	if _, err := r.builder.
		Update(`"user"`).
		Set(`"email"`, m.Email).
		Set(`"first_name"`, m.FirstName).
		Set(`"last_name"`, m.LastName).
		Set(`"refresh_token_hash"`, m.RefreshTokenHash).
		Set(`"google_id"`, m.GoogleID).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) IsEmailExists(ctx context.Context, email string) (bool, error) {
	if _, err := r.Get(ctx, user.ByEmail(email)); err != nil {
		if errdefs.IsUserNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *RepositoryCE) GetOrganizationAccess(ctx context.Context, opts ...user.OrganizationAccessRepositoryOption) (*user.UserOrganizationAccess, error) {
	query, args, err := r.buildOrganizationAccessQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := user.UserOrganizationAccess{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserOrganizationAccessNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) ListOrganizationAccesses(ctx context.Context, opts ...user.OrganizationAccessRepositoryOption) ([]*user.UserOrganizationAccess, error) {
	query, args, err := r.buildOrganizationAccessQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*user.UserOrganizationAccess, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *RepositoryCE) buildOrganizationAccessQuery(ctx context.Context, opts ...user.OrganizationAccessRepositoryOption) (string, []any, error) {
	q := r.builder.Select(
		`uoa."id"`,
		`uoa."user_id"`,
		`uoa."organization_id"`,
		`uoa."role"`,
		`uoa."created_at"`,
		`uoa."updated_at"`,
	).
		From(`"user_organization_access" uoa`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (r *RepositoryCE) CreateOrganizationAccess(ctx context.Context, m *user.UserOrganizationAccess) error {
	if _, err := r.builder.
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
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) UpdateOrganizationAccess(ctx context.Context, m *user.UserOrganizationAccess) error {
	if _, err := r.builder.
		Update(`"user_organization_access"`).
		Set(`"role"`, m.Role).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) DeleteOrganizationAccess(ctx context.Context, m *user.UserOrganizationAccess) error {
	if _, err := r.builder.
		Delete(`"user_organization_access"`).
		Where(sq.Eq{`"user_id"`: m.UserID, `"organization_id"`: m.OrganizationID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		if err == sql.ErrNoRows {
			return errdefs.ErrUserOrganizationAccessNotFound(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) GetGroup(ctx context.Context, opts ...user.GroupRepositoryOption) (*user.UserGroup, error) {
	query, args, err := r.buildGroupQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := user.UserGroup{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserGroupNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) ListGroups(ctx context.Context, opts ...user.GroupRepositoryOption) ([]*user.UserGroup, error) {
	query, args, err := r.buildGroupQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*user.UserGroup, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *RepositoryCE) buildGroupQuery(ctx context.Context, opts ...user.GroupRepositoryOption) (string, []any, error) {
	q := r.builder.Select(
		`ug."id"`,
		`ug."user_id"`,
		`ug."group_id"`,
		`ug."created_at"`,
		`ug."updated_at"`,
	).
		From(`"user_group" ug`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (r *RepositoryCE) BulkInsertGroups(ctx context.Context, m []*user.UserGroup) error {
	if len(m) == 0 {
		return nil
	}

	q := r.builder.
		Insert(`"user_group"`).
		Columns(`"id"`, `"user_id"`, `"group_id"`)

	for _, v := range m {
		q = q.Values(v.ID, v.UserID, v.GroupID)
	}

	if _, err := q.
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) BulkDeleteGroups(ctx context.Context, m []*user.UserGroup) error {
	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *user.UserGroup, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := r.builder.
		Delete(`"user_group"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) GetInvitation(ctx context.Context, opts ...user.InvitationRepositoryOption) (*user.UserInvitation, error) {
	query, args, err := r.buildInvitationQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := user.UserInvitation{}
	if err := r.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserInvitationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (r *RepositoryCE) ListInvitations(ctx context.Context, opts ...user.InvitationRepositoryOption) ([]*user.UserInvitation, error) {
	query, args, err := r.buildInvitationQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*user.UserInvitation, 0)
	if err := r.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (r *RepositoryCE) buildInvitationQuery(ctx context.Context, opts ...user.InvitationRepositoryOption) (string, []any, error) {
	q := r.builder.Select(
		`ui."id"`,
		`ui."organization_id"`,
		`ui."email"`,
		`ui."role"`,
		`ui."created_at"`,
		`ui."updated_at"`,
	).
		From(`"user_invitation" ui`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (r *RepositoryCE) DeleteInvitation(ctx context.Context, m *user.UserInvitation) error {
	if _, err := r.builder.
		Delete(`"user_invitation"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) BulkInsertInvitations(ctx context.Context, m []*user.UserInvitation) error {
	if len(m) == 0 {
		return nil
	}

	q := r.builder.
		Insert(`"user_invitation"`).
		Columns(`"id"`, `"organization_id"`, `"email"`, `"role"`)

	for _, v := range m {
		q = q.Values(v.ID, v.OrganizationID, v.Email, v.Role)
	}

	if _, err := q.
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (r *RepositoryCE) IsInvitationEmailExists(ctx context.Context, orgID uuid.UUID, email string) (bool, error) {
	if _, err := r.GetInvitation(ctx, user.InvitationByOrganizationID(orgID), user.InvitationByEmail(email)); err != nil {
		if errdefs.IsUserInvitationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
