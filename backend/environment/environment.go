package environment

import (
	"crypto/rand"
	"fmt"
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

func (e *Environment) GenerateAPIKey() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomStr := make([]byte, 51)
	for i := range randomStr {
		randomStr[i] = charset[randomBytes[i%len(randomBytes)]%byte(len(charset))]
	}

	return fmt.Sprintf("%s_%s", e.Slug, string(randomStr)), nil
}
