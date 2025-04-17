package user

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/config"
)

const (
	EmailTokenExpiration     = time.Duration(24) * time.Hour
	tokenExpiration          = time.Duration(60) * time.Minute
	tokenExpirationDev       = time.Duration(365*24) * time.Hour
	RefreshTokenExpiration   = time.Duration(30*24) * time.Hour
	XSRFTokenExpiration      = time.Duration(30*24) * time.Hour
	RefreshTokenMaxAgeBuffer = time.Duration(7*24) * time.Hour
	TmpTokenExpiration       = time.Duration(30) * time.Minute

	SaveAuthPath = "/api/v1/auth/save"
)

func TokenExpiration() time.Duration {
	if config.Config.Env == config.EnvLocal {
		return tokenExpirationDev
	}
	return tokenExpiration
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
