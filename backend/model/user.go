package model

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

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
	Secret               string     `db:"secret"`
	GoogleID             string     `db:"google_id"`
	EmailAuthenticatedAt *time.Time `db:"email_authenticated_at"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
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

type UserStore interface {
	Get(context.Context, ...storeopts.UserOption) (*User, error)
	List(context.Context, ...storeopts.UserOption) ([]*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	IsEmailExists(context.Context, string) (bool, error)

	GetOrganizationAccess(context.Context, ...storeopts.UserOrganizationAccessOption) (*UserOrganizationAccess, error)
	ListOrganizationAccesses(context.Context, ...storeopts.UserOrganizationAccessOption) ([]*UserOrganizationAccess, error)
	CreateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	UpdateOrganizationAccess(context.Context, *UserOrganizationAccess) error
	DeleteOrganizationAccess(context.Context, *UserOrganizationAccess) error

	GetGroup(context.Context, ...storeopts.UserGroupOption) (*UserGroup, error)
	ListGroups(context.Context, ...storeopts.UserGroupOption) ([]*UserGroup, error)
	BulkInsertGroups(context.Context, []*UserGroup) error
	BulkDeleteGroups(context.Context, []*UserGroup) error

	GetInvitation(context.Context, ...storeopts.UserInvitationOption) (*UserInvitation, error)
	ListInvitations(context.Context, ...storeopts.UserInvitationOption) ([]*UserInvitation, error)
	DeleteInvitation(context.Context, *UserInvitation) error
	BulkInsertInvitations(context.Context, []*UserInvitation) error
	IsInvitationEmailExists(context.Context, uuid.UUID, string) (bool, error)
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

// Email structure for sending magic link email.
type SendMagicLinkEmail struct {
	To        string
	FirstName string
	URL       string
}

// Email structure for sending multiple organizations magic link email.
type SendMultipleOrganizationsMagicLinkEmail struct {
	To        string
	FirstName string
	Email     string
	LoginURLs []string
}

// Email structure for sending multiple organizations login email.
type SendMultipleOrganizationsLoginEmail struct {
	To        string
	FirstName string
	Email     string
	LoginURLs []string
}

// SendInvitationMagicLinkEmail represents the data needed to send an invitation magic link email.
type SendInvitationMagicLinkEmail struct {
	To        string
	URL       string
	FirstName string
}

type UserMailer interface {
	SendUpdateEmailInstructions(ctx context.Context, in *SendUpdateUserEmailInstructions) error
	SendInvitationEmail(ctx context.Context, in *SendInvitationEmail) error
	SendMagicLinkEmail(ctx context.Context, in *SendMagicLinkEmail) error
	SendMultipleOrganizationsMagicLinkEmail(ctx context.Context, in *SendMultipleOrganizationsMagicLinkEmail) error
	SendMultipleOrganizationsLoginEmail(ctx context.Context, in *SendMultipleOrganizationsLoginEmail) error
	SendInvitationMagicLinkEmail(ctx context.Context, in *SendInvitationMagicLinkEmail) error
}
