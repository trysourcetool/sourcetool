package core

import (
	"crypto/sha256"
	"encoding/hex"
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

func HashAPIKey(plainAPIKey string) string {
	hash := sha256.Sum256([]byte(plainAPIKey))
	return hex.EncodeToString(hash[:])
}
