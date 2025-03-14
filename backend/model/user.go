package model

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

const (
	EmailTokenExpiration = time.Duration(24) * time.Hour
	tokenExpiration      = time.Duration(60) * time.Minute
	tokenExpirationDev   = time.Duration(365*24) * time.Hour
	SecretExpiration     = time.Duration(30*24) * time.Hour
	XSRFTokenExpiration  = time.Duration(30*24) * time.Hour
	SecretMaxAgeBuffer   = time.Duration(7*24) * time.Hour
	TmpTokenExpiration   = time.Duration(30) * time.Minute

	SaveAuthPath = "/api/v1/users/saveAuth"
)

func TokenExpiration() time.Duration {
	if config.Config.Env == config.EnvLocal {
		return tokenExpirationDev
	}
	return tokenExpiration
}

type User struct {
	ID                   uuid.UUID  `db:"id"`
	Email                string     `db:"email"`
	FirstName            string     `db:"first_name"`
	LastName             string     `db:"last_name"`
	Password             string     `db:"password"`
	Secret               string     `db:"secret"`
	GoogleID             string     `db:"google_id"`
	EmailAuthenticatedAt *time.Time `db:"email_authenticated_at"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
}

type UserRegistrationRequest struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserGoogleAuthRequest struct {
	ID        uuid.UUID `db:"id"`
	GoogleID  string    `db:"google_id"`
	Email     string    `db:"email"`
	Domain    string    `db:"domain"`
	Invited   bool      `db:"invited"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserOrganizationRole int

const (
	UserOrganizationRoleUnknown UserOrganizationRole = iota
	UserOrganizationRoleAdmin
	UserOrganizationRoleDeveloper
	UserOrganizationRoleMember

	userOrganizationRoleUnknown   = "unknown"
	userOrganizationRoleAdmin     = "admin"
	userOrganizationRoleDeveloper = "developer"
	userOrganizationRoleMember    = "member"
)

func (r UserOrganizationRole) String() string {
	switch r {
	case UserOrganizationRoleAdmin:
		return userOrganizationRoleAdmin
	case UserOrganizationRoleDeveloper:
		return userOrganizationRoleDeveloper
	case UserOrganizationRoleMember:
		return userOrganizationRoleMember
	}
	return userOrganizationRoleUnknown
}

func UserOrganizationRoleFromString(s string) UserOrganizationRole {
	switch s {
	case userOrganizationRoleAdmin:
		return UserOrganizationRoleAdmin
	case userOrganizationRoleDeveloper:
		return UserOrganizationRoleDeveloper
	case userOrganizationRoleMember:
		return UserOrganizationRoleMember
	}
	return UserOrganizationRoleUnknown
}

type UserInvitation struct {
	ID             uuid.UUID            `db:"id"`
	OrganizationID uuid.UUID            `db:"organization_id"`
	Email          string               `db:"email"`
	Role           UserOrganizationRole `db:"role"`
	CreatedAt      time.Time            `db:"created_at"`
	UpdatedAt      time.Time            `db:"updated_at"`
}

type UserOrganizationAccess struct {
	ID             uuid.UUID            `db:"id"`
	UserID         uuid.UUID            `db:"user_id"`
	OrganizationID uuid.UUID            `db:"organization_id"`
	Role           UserOrganizationRole `db:"role"`
	CreatedAt      time.Time            `db:"created_at"`
	UpdatedAt      time.Time            `db:"updated_at"`
}

type UserGroup struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	GroupID   uuid.UUID `db:"group_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *User) FullName() string {
	return fmt.Sprintf("%s %s", m.FirstName, m.LastName)
}

func ValidatePassword(p string) error {
	if err := validation.Validate(p, validation.Length(8, 32)); err != nil {
		return err
	}

	return nil
}

type UserStore interface {
	Get(context.Context, ...storeopts.UserOption) (*User, error)
	List(context.Context, ...storeopts.UserOption) ([]*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	IsEmailExists(context.Context, string) (bool, error)

	GetRegistrationRequest(context.Context, ...storeopts.UserRegistrationRequestOption) (*UserRegistrationRequest, error)
	CreateRegistrationRequest(context.Context, *UserRegistrationRequest) error
	DeleteRegistrationRequest(context.Context, *UserRegistrationRequest) error
	IsRegistrationRequestExists(context.Context, string) (bool, error)

	GetOrganizationAccess(context.Context, ...storeopts.UserOrganizationAccessOption) (*UserOrganizationAccess, error)
	ListOrganizationAccesses(context.Context, ...storeopts.UserOrganizationAccessOption) ([]*UserOrganizationAccess, error)
	CreateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	UpdateOrganizationAccess(context.Context, *UserOrganizationAccess) error

	GetGroup(context.Context, ...storeopts.UserGroupOption) (*UserGroup, error)
	ListGroups(context.Context, ...storeopts.UserGroupOption) ([]*UserGroup, error)
	BulkInsertGroups(context.Context, []*UserGroup) error
	BulkDeleteGroups(context.Context, []*UserGroup) error

	GetInvitation(context.Context, ...storeopts.UserInvitationOption) (*UserInvitation, error)
	ListInvitations(context.Context, ...storeopts.UserInvitationOption) ([]*UserInvitation, error)
	DeleteInvitation(context.Context, *UserInvitation) error
	BulkInsertInvitations(context.Context, []*UserInvitation) error
	IsInvitationEmailExists(context.Context, uuid.UUID, string) (bool, error)

	GetGoogleAuthRequest(context.Context, uuid.UUID) (*UserGoogleAuthRequest, error)
	ListExpiredGoogleAuthRequests(context.Context) ([]*UserGoogleAuthRequest, error)
	CreateGoogleAuthRequest(context.Context, *UserGoogleAuthRequest) error
	UpdateGoogleAuthRequest(context.Context, *UserGoogleAuthRequest) error
	DeleteGoogleAuthRequest(context.Context, *UserGoogleAuthRequest) error
	BulkDeleteGoogleAuthRequests(context.Context, []*UserGoogleAuthRequest) error
}

const (
	UserSignatureSubjectEmail             = "email"
	UserSignatureSubjectUpdateEmail       = "update_email"
	UserSignatureSubjectActivate          = "activate"
	UserSignatureSubjectInvitation        = "invitaiton"
	UserSignatureSubjectGoogleAuthRequest = "google_auth_request"
)

type UserClaims struct {
	UserID string
	Email  string
	jwt.RegisteredClaims
}

type UserEmailClaims struct {
	Email string
	jwt.RegisteredClaims
}

type UserGoogleAuthRequestClaims struct {
	GoogleAuthRequestID string
	jwt.RegisteredClaims
}

type UserAuthClaims struct {
	UserID string
	// OrganizationID string
	Email     string
	XSRFToken string
	jwt.RegisteredClaims
}

type UserSigner interface {
	SignedString(context.Context, *UserClaims) (string, error)
	SignedStringFromEmail(context.Context, *UserEmailClaims) (string, error)
	SignedStringGoogleAuthRequest(context.Context, *UserGoogleAuthRequestClaims) (string, error)
	SignedStringAuth(context.Context, *UserAuthClaims) (string, error)
	ClaimsFromToken(context.Context, string) (*UserClaims, error)
	EmailClaimsFromToken(context.Context, string) (*UserEmailClaims, error)
	GoogleAuthRequestClaimsFromToken(context.Context, string) (*UserGoogleAuthRequestClaims, error)
	AuthClaimsFromToken(context.Context, string) (*UserAuthClaims, error)
}

type SendSignUpInstructions struct {
	To  string
	URL string
}

type SendUpdateUserEmailInstructions struct {
	To        string
	FirstName string
	URL       string
}

type SendWelcomeEmail struct {
	To string
}

type SendInvitationEmail struct {
	Invitees string
	URLs     map[string]string // email -> url
}

type UserMailer interface {
	SendSignUpInstructions(context.Context, *SendSignUpInstructions) error
	SendUpdateEmailInstructions(context.Context, *SendUpdateUserEmailInstructions) error
	SendInvitationEmail(context.Context, *SendInvitationEmail) error
}
