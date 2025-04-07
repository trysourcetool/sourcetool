package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// hashSecret creates a SHA-256 hash of a plaintext secret.
func hashSecret(plainSecret string) string {
	hash := sha256.Sum256([]byte(plainSecret))
	return hex.EncodeToString(hash[:])
}

// generateSecret creates a new random secret and its hash for user authentication.
func generateSecret() (plainSecret, hashedSecret string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	plainSecret = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
	hashedSecret = hashSecret(plainSecret)

	return plainSecret, hashedSecret, nil
}
