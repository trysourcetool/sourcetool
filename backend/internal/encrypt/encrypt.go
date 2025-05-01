package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

type Encryptor struct {
	gcm cipher.AEAD
}

func NewEncryptor() (*Encryptor, error) {
	keyB64 := config.Config.EncryptionKey
	if keyB64 == "" {
		return nil, errors.New("ENCRYPTION_KEY not set")
	}
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil || len(key) != 32 {
		return nil, errors.New("key must be 32byte base64")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Encryptor{gcm: gcm}, nil
}

func (e *Encryptor) Encrypt(plain []byte) (nonce, cipherText []byte, err error) {
	nonce = make([]byte, e.gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}
	cipherText = e.gcm.Seal(nil, nonce, plain, nil)
	return
}

func (e *Encryptor) Decrypt(nonce, cipherText []byte) ([]byte, error) {
	return e.gcm.Open(nil, nonce, cipherText, nil)
}
