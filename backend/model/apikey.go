package model

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool/backend/storeopts"
)

type APIKey struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	EnvironmentID  uuid.UUID `db:"environment_id"`
	UserID         uuid.UUID `db:"user_id"`
	Name           string    `db:"name"`
	Key            string    `db:"key"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func GenerateAPIKey(orgSubdomain, envSlug string) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	length := 60
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	result := []byte{}
	for _, b := range randomBytes {
		result = append(result, charset[int(b)%len(charset)])
	}

	key := fmt.Sprintf("%s_%s_%s", orgSubdomain, envSlug, string(result))

	return key[:length], nil
}

type APIKeyStore interface {
	Get(context.Context, ...storeopts.APIKeyOption) (*APIKey, error)
	List(context.Context, ...storeopts.APIKeyOption) ([]*APIKey, error)
	Create(context.Context, *APIKey) error
	Update(context.Context, *APIKey) error
	Delete(context.Context, *APIKey) error
}
