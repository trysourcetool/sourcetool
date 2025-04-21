package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/sourcetool/backend/internal/core"
	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func (db *DB) GetSession(ctx context.Context, queries ...SessionQuery) (*core.Session, error) {
	query, args, err := db.buildSessionQuery(ctx, queries...)
	if err != nil {
		return nil, err
	}

	m := core.Session{}
	if err := db.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrSessionNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (db *DB) buildSessionQuery(ctx context.Context, queries ...SessionQuery) (string, []any, error) {
	q := db.builder.Select(
		`s."id"`,
		`s."organization_id"`,
		`s."user_id"`,
		`s."api_key_id"`,
		`s."host_instance_id"`,
		`s."created_at"`,
		`s."updated_at"`,
	).
		From(`"session" s`)

	for _, query := range queries {
		q = query.apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (db *DB) CreateSession(ctx context.Context, tx *sqlx.Tx, m *core.Session) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Insert(`"session"`).
		Columns(
			`"id"`,
			`"organization_id"`,
			`"user_id"`,
			`"api_key_id"`,
			`"host_instance_id"`,
		).
		Values(
			m.ID,
			m.OrganizationID,
			m.UserID,
			m.APIKeyID,
			m.HostInstanceID,
		).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (db *DB) DeleteSession(ctx context.Context, tx *sqlx.Tx, m *core.Session) error {
	var runner sq.BaseRunner
	if tx != nil {
		runner = tx
	} else {
		runner = db.db
	}

	if _, err := db.builder.
		Delete(`"session"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(runner).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
