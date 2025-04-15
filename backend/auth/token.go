package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// hashRefreshToken creates a SHA-256 hash of a plaintext refresh token.
func hashRefreshToken(plainRefreshToken string) string {
	hash := sha256.Sum256([]byte(plainRefreshToken))
	return hex.EncodeToString(hash[:])
}

// generateRefreshToken creates a new random refresh token and its hash for user authentication.
func generateRefreshToken() (plainRefreshToken, hashedRefreshToken string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	plainRefreshToken = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
	hashedRefreshToken = hashRefreshToken(plainRefreshToken)

	return plainRefreshToken, hashedRefreshToken, nil
}
