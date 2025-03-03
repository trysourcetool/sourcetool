package user

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/model"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

type StoreCE struct {
	db      infra.DB
	builder sq.StatementBuilderType
}

func NewStoreCE(db infra.DB) *StoreCE {
	return &StoreCE{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *StoreCE) Get(ctx context.Context, opts ...storeopts.UserOption) (*model.User, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.User{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) List(ctx context.Context, opts ...storeopts.UserOption) ([]*model.User, error) {
	query, args, err := s.buildQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.User, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) buildQuery(ctx context.Context, opts ...storeopts.UserOption) (string, []any, error) {
	q := s.builder.Select(
		`u."id"`,
		`u."created_at"`,
		`u."email"`,
		`u."first_name"`,
		`u."last_name"`,
		`u."updated_at"`,
		`u."password"`,
		`u."secret"`,
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

func (s *StoreCE) Create(ctx context.Context, m *model.User) error {
	if _, err := s.builder.
		Insert(`"user"`).
		Columns(
			`"id"`,
			`"email"`,
			`"first_name"`,
			`"last_name"`,
			`"password"`,
			`"secret"`,
			`"google_id"`,
		).
		Values(
			m.ID,
			m.Email,
			m.FirstName,
			m.LastName,
			m.Password,
			m.Secret,
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

func (s *StoreCE) Update(ctx context.Context, m *model.User) error {
	if _, err := s.builder.
		Update(`"user"`).
		Set(`"email"`, m.Email).
		Set(`"first_name"`, m.FirstName).
		Set(`"last_name"`, m.LastName).
		Set(`"password"`, m.Password).
		Set(`"secret"`, m.Secret).
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

func (s *StoreCE) IsEmailExists(ctx context.Context, email string) (bool, error) {
	if _, err := s.Get(ctx, storeopts.UserByEmail(email)); err != nil {
		if errdefs.IsUserNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *StoreCE) GetRegistrationRequest(ctx context.Context, opts ...storeopts.UserRegistrationRequestOption) (*model.UserRegistrationRequest, error) {
	query, args, err := s.buildRegistrationRequestQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.UserRegistrationRequest{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserRegistrationRequestNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) buildRegistrationRequestQuery(ctx context.Context, opts ...storeopts.UserRegistrationRequestOption) (string, []any, error) {
	q := s.builder.Select(
		`urr."id"`,
		`urr."email"`,
		`urr."created_at"`,
		`urr."updated_at"`,
	).
		From(`"user_registration_request" urr`)

	for _, o := range opts {
		q = o.Apply(q)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return "", nil, errdefs.ErrDatabase(err)
	}

	return query, args, err
}

func (s *StoreCE) CreateRegistrationRequest(ctx context.Context, m *model.UserRegistrationRequest) error {
	if _, err := s.builder.
		Insert(`"user_registration_request"`).
		Columns(`"id"`, `"email"`).
		Values(m.ID, m.Email).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) DeleteRegistrationRequest(ctx context.Context, m *model.UserRegistrationRequest) error {
	if _, err := s.builder.
		Delete(`"user_registration_request"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) IsRegistrationRequestExists(ctx context.Context, email string) (bool, error) {
	if _, err := s.GetRegistrationRequest(ctx, storeopts.UserRegistrationRequestByEmail(email)); err != nil {
		if errdefs.IsUserRegistrationRequestNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *StoreCE) GetOrganizationAccess(ctx context.Context, opts ...storeopts.UserOrganizationAccessOption) (*model.UserOrganizationAccess, error) {
	query, args, err := s.buildOrganizationAccessQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.UserOrganizationAccess{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserOrganizationAccessNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) ListOrganizationAccesses(ctx context.Context, opts ...storeopts.UserOrganizationAccessOption) ([]*model.UserOrganizationAccess, error) {
	query, args, err := s.buildOrganizationAccessQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.UserOrganizationAccess, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) buildOrganizationAccessQuery(ctx context.Context, opts ...storeopts.UserOrganizationAccessOption) (string, []any, error) {
	q := s.builder.Select(
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

func (s *StoreCE) CreateOrganizationAccess(ctx context.Context, m *model.UserOrganizationAccess) error {
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

func (s *StoreCE) UpdateOrganizationAccess(ctx context.Context, m *model.UserOrganizationAccess) error {
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

func (s *StoreCE) GetGroup(ctx context.Context, opts ...storeopts.UserGroupOption) (*model.UserGroup, error) {
	query, args, err := s.buildGroupQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.UserGroup{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserGroupNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) ListGroups(ctx context.Context, opts ...storeopts.UserGroupOption) ([]*model.UserGroup, error) {
	query, args, err := s.buildGroupQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.UserGroup, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) buildGroupQuery(ctx context.Context, opts ...storeopts.UserGroupOption) (string, []any, error) {
	q := s.builder.Select(
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

func (s *StoreCE) BulkInsertGroups(ctx context.Context, m []*model.UserGroup) error {
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

func (s *StoreCE) BulkDeleteGroups(ctx context.Context, m []*model.UserGroup) error {
	if len(m) == 0 {
		return nil
	}

	ids := lo.Map(m, func(x *model.UserGroup, _ int) uuid.UUID {
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

func (s *StoreCE) GetInvitation(ctx context.Context, opts ...storeopts.UserInvitationOption) (*model.UserInvitation, error) {
	query, args, err := s.buildInvitationQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := model.UserInvitation{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserInvitationNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) ListInvitations(ctx context.Context, opts ...storeopts.UserInvitationOption) ([]*model.UserInvitation, error) {
	query, args, err := s.buildInvitationQuery(ctx, opts...)
	if err != nil {
		return nil, err
	}

	m := make([]*model.UserInvitation, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) buildInvitationQuery(ctx context.Context, opts ...storeopts.UserInvitationOption) (string, []any, error) {
	q := s.builder.Select(
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

func (s *StoreCE) DeleteInvitation(ctx context.Context, m *model.UserInvitation) error {
	if _, err := s.builder.
		Delete(`"user_invitation"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) BulkInsertInvitations(ctx context.Context, m []*model.UserInvitation) error {
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

func (s *StoreCE) IsInvitationEmailExists(ctx context.Context, orgID uuid.UUID, email string) (bool, error) {
	if _, err := s.GetInvitation(ctx, storeopts.UserInvitationByOrganizationID(orgID), storeopts.UserInvitationByEmail(email)); err != nil {
		if errdefs.IsUserInvitationNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *StoreCE) GetGoogleAuthRequest(ctx context.Context, state uuid.UUID) (*model.UserGoogleAuthRequest, error) {
	query, args, err := s.builder.
		Select(`"id"`, `"google_id"`, `"email"`, `"domain"`, `"expires_at"`, `"invited"`, `"created_at"`, `"updated_at"`).
		From(`"user_google_auth_request"`).
		Where(sq.Eq{`"id"`: state}).
		ToSql()
	if err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	m := model.UserGoogleAuthRequest{}
	if err := s.db.GetContext(ctx, &m, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserGoogleAuthRequestNotFound(err)
		}
		return nil, errdefs.ErrDatabase(err)
	}

	return &m, nil
}

func (s *StoreCE) ListExpiredGoogleAuthRequests(ctx context.Context) ([]*model.UserGoogleAuthRequest, error) {
	query, args, err := s.builder.
		Select(`"id"`, `"google_id"`, `"email"`, `"domain"`, `"expires_at"`, `"invited"`, `"created_at"`, `"updated_at"`).
		From(`"user_google_auth_request"`).
		Where(`"expires_at" < $1`, time.Now()).
		ToSql()
	if err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	m := make([]*model.UserGoogleAuthRequest, 0)
	if err := s.db.SelectContext(ctx, &m, query, args...); err != nil {
		return nil, errdefs.ErrDatabase(err)
	}

	return m, nil
}

func (s *StoreCE) CreateGoogleAuthRequest(ctx context.Context, m *model.UserGoogleAuthRequest) error {
	if _, err := s.builder.
		Insert(`"user_google_auth_request"`).
		Columns(`"id"`, `"google_id"`, `"email"`, `"domain"`, `"expires_at"`, `"invited"`).
		Values(m.ID, m.GoogleID, m.Email, m.Domain, m.ExpiresAt, m.Invited).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) UpdateGoogleAuthRequest(ctx context.Context, m *model.UserGoogleAuthRequest) error {
	if _, err := s.builder.
		Update(`"user_google_auth_request"`).
		Set(`"google_id"`, m.GoogleID).
		Set(`"email"`, m.Email).
		Set(`"domain"`, m.Domain).
		Set(`"expires_at"`, m.ExpiresAt).
		Set(`"invited"`, m.Invited).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) DeleteGoogleAuthRequest(ctx context.Context, m *model.UserGoogleAuthRequest) error {
	if _, err := s.builder.
		Delete(`"user_google_auth_request"`).
		Where(sq.Eq{`"id"`: m.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *StoreCE) BulkDeleteGoogleAuthRequests(ctx context.Context, m []*model.UserGoogleAuthRequest) error {
	ids := lo.Map(m, func(x *model.UserGoogleAuthRequest, _ int) uuid.UUID {
		return x.ID
	})

	if _, err := s.builder.
		Delete(`"user_google_auth_request"`).
		Where(sq.Eq{`"id"`: ids}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}
