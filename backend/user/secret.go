package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"

	"github.com/trysourcetool/sourcetool/backend/errdefs"
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

// hashPassword creates a bcrypt hash of the password.
func hashPassword(password string) (string, error) {
	encodedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errdefs.ErrInternal(err)
	}

	return hex.EncodeToString(encodedPass[:]), nil
}

// verifyPassword checks if the provided password matches the stored hash.
func verifyPassword(storedHash, password string) error {
	h, err := hex.DecodeString(storedHash)
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	if err = bcrypt.CompareHashAndPassword(h, []byte(password)); err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	return nil
}
