package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func (db *DB) GetUser(ctx context.Context, queries ...UserQuery) (*core.User, error) {
	query, args, err := db.buildUserQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.User{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) ListUsers(ctx context.Context, queries ...UserQuery) ([]*core.User, error) {
	query, args, err := db.buildUserQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.User, 0)
	if err := db.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (db *DB) buildUserQuery(ctx context.Context, queries ...UserQuery) (string, []any, error) {
	q := db.builder.Select(
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

	q = applyUserQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) CreateUser(ctx context.Context, tx *sqlx.Tx, m *core.User) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
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
		RunWith(runner).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) UpdateUser(ctx context.Context, tx *sqlx.Tx, m *core.User) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Update(`"user"`).
		Set(`"email"`, m.Email).
		Set(`"first_name"`, m.FirstName).
		Set(`"last_name"`, m.LastName).
		Set(`"refresh_token_hash"`, m.RefreshTokenHash).
		Set(`"google_id"`, m.GoogleID).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) IsUserEmailExists(ctx context.Context, email string) (bool, error) {
	if _, err := db.GetUser(ctx, UserByEmail(email)); err != nil {
		if errdefs.IsUserNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *DB) GetUserOrganizationAccess(ctx context.Context, queries ...UserOrganizationAccessQuery) (*core.UserOrganizationAccess, error) {
	query, args, err := db.buildUserOrganizationAccessQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.UserOrganizationAccess{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserOrganizationAccessNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) ListUserOrganizationAccesses(ctx context.Context, queries ...UserOrganizationAccessQuery) ([]*core.UserOrganizationAccess, error) {
	query, args, err := db.buildUserOrganizationAccessQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.UserOrganizationAccess, 0)
	if err := db.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (db *DB) buildUserOrganizationAccessQuery(ctx context.Context, queries ...UserOrganizationAccessQuery) (string, []any, error) {
	q := db.builder.Select(
		`uoa."id"`,
		`uoa."user_id"`,
		`uoa."organization_id"`,
		`uoa."role"`,
		`uoa."created_at"`,
		`uoa."updated_at"`,
	).
		From(`"user_organization_access" uoa`)

	q = applyUserOrganizationAccessQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) CreateUserOrganizationAccess(ctx context.Context, tx *sqlx.Tx, m *core.UserOrganizationAccess) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
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
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) UpdateUserOrganizationAccess(ctx context.Context, tx *sqlx.Tx, m *core.UserOrganizationAccess) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Update(`"user_organization_access"`).
		Set(`"role"`, m.Role).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) DeleteUserOrganizationAccess(ctx context.Context, tx *sqlx.Tx, m *core.UserOrganizationAccess) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Delete(`"user_organization_access"`).
		Where(sq.Eq{`"user_id"`: m.UserID, `"organization_id"`: m.OrganizationID}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		if err == sql.ErrNoRows {
			return errdefs.ErrUserOrganizationAccessNotFound(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) GetUserGroup(ctx context.Context, queries ...UserGroupQuery) (*core.UserGroup, error) {
	query, args, err := db.buildUserGroupQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.UserGroup{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserGroupNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) ListUserGroups(ctx context.Context, queries ...UserGroupQuery) ([]*core.UserGroup, error) {
	query, args, err := db.buildUserGroupQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.UserGroup, 0)
	if err := db.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (db *DB) buildUserGroupQuery(ctx context.Context, queries ...UserGroupQuery) (string, []any, error) {
	q := db.builder.Select(
		`ug."id"`,
		`ug."user_id"`,
		`ug."group_id"`,
		`ug."created_at"`,
		`ug."updated_at"`,
	).
		From(`"user_group" ug`)

	q = applyUserGroupQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) BulkInsertUserGroups(ctx context.Context, tx *sqlx.Tx, m []*core.UserGroup) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(m) == 0 {
		return nil
	}

	q := db.builder.
		Insert(`"user_group"`).
		Columns(`"id"`, `"user_id"`, `"group_id"`)

	for _, v := range m {
		q = q.Values(v.ID, v.UserID, v.GroupID)
	}

	if _, err := q.
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) BulkDeleteUserGroups(ctx context.Context, tx *sqlx.Tx, m []*core.UserGroup) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *core.UserGroup, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := db.builder.
		Delete(`"user_group"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) GetUserInvitation(ctx context.Context, queries ...UserInvitationQuery) (*core.UserInvitation, error) {
	query, args, err := db.buildUserInvitationQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.UserInvitation{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserInvitationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) ListUserInvitations(ctx context.Context, queries ...UserInvitationQuery) ([]*core.UserInvitation, error) {
	query, args, err := db.buildUserInvitationQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := make([]*core.UserInvitation, 0)
	if err := db.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (db *DB) buildUserInvitationQuery(ctx context.Context, queries ...UserInvitationQuery) (string, []any, error) {
	q := db.builder.Select(
		`ui."id"`,
		`ui."organization_id"`,
		`ui."email"`,
		`ui."role"`,
		`ui."created_at"`,
		`ui."updated_at"`,
	).
		From(`"user_invitation" ui`)

	q = applyUserInvitationQueries(q, queries...)

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) DeleteUserInvitation(ctx context.Context, tx *sqlx.Tx, m *core.UserInvitation) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Delete(`"user_invitation"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) BulkInsertUserInvitations(ctx context.Context, tx *sqlx.Tx, m []*core.UserInvitation) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if len(m) == 0 {
		return nil
	}

	q := db.builder.
		Insert(`"user_invitation"`).
		Columns(`"id"`, `"organization_id"`, `"email"`, `"role"`)

	for _, v := range m {
		q = q.Values(v.ID, v.OrganizationID, v.Email, v.Role)
	}

	if _, err := q.
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) IsUserInvitationEmailExists(ctx context.Context, orgID uuid.UUID, email string) (bool, error) {
	if _, err := db.GetUserInvitation(ctx, UserInvitationByOrganizationID(orgID), UserInvitationByEmail(email)); err != nil {
		if errdefs.IsUserInvitationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
