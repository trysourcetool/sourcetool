package core

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/encrypt"
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

func (e *Environment) GenerateAPIKey() (plainAPIKey, hashedAPIKey string, ciphertext, nonce []byte, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", nil, nil, err
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomStr := make([]byte, 51)
	for i := range randomStr {
		randomStr[i] = charset[randomBytes[i%len(randomBytes)]%byte(len(charset))]
	}

	plainAPIKey = fmt.Sprintf("%s_%s", e.Slug, string(randomStr))
	hashedAPIKey = HashAPIKey(plainAPIKey)

	encryptor, err := encrypt.NewEncryptor()
	if err != nil {
		return "", "", nil, nil, err
	}

	ciphertext, nonce, err = encryptor.Encrypt([]byte(plainAPIKey))
	if err != nil {
		return "", "", nil, nil, err
	}

	return plainAPIKey, hashedAPIKey, ciphertext, nonce, nil
}
