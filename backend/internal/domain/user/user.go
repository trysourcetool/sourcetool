package user

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

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

type User struct {
	ID               uuid.UUID `db:"id"`
	Email            string    `db:"email"`
	FirstName        string    `db:"first_name"`
	LastName         string    `db:"last_name"`
	RefreshTokenHash string    `db:"refresh_token_hash"`
	GoogleID         string    `db:"google_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
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
