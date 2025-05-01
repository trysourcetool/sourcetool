package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type APIKey struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	EnvironmentID  uuid.UUID `db:"environment_id"`
	UserID         uuid.UUID `db:"user_id"`
	Name           string    `db:"name"`
	KeyHash        string    `db:"key_hash"`
	KeyCiphertext  []byte    `db:"key_ciphertext"`
	KeyNonce       []byte    `db:"key_nonce"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func GenerateAPIKey(slug string) (plainAPIKey string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomStr := make([]byte, 51)
	for i := range randomStr {
		randomStr[i] = charset[randomBytes[i%len(randomBytes)]%byte(len(charset))]
	}

	plainAPIKey = fmt.Sprintf("%s_%s", slug, string(randomStr))

	return plainAPIKey, nil
}

func HashAPIKey(plainAPIKey string) string {
	hash := sha256.Sum256([]byte(plainAPIKey))
	return hex.EncodeToString(hash[:])
}
