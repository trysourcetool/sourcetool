package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	EnvironmentNameProduction   = "Production"
	EnvironmentNameDevelopment  = "Development"
	EnvironmentSlugProduction   = "production"
	EnvironmentSlugDevelopment  = "development"
	EnvironmentColorProduction  = "#9333EA"
	EnvironmentColorDevelopment = "#33ADEA"
)

type Environment struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Name           string    `db:"name"`
	Slug           string    `db:"slug"`
	Color          string    `db:"color"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
