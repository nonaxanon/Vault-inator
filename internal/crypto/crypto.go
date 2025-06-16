package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

// KeyDerivationParams holds the parameters for Argon2id key derivation.
type KeyDerivationParams struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultKeyDerivationParams returns recommended parameters for Argon2id.
func DefaultKeyDerivationParams() KeyDerivationParams {
	return KeyDerivationParams{
		Time:    1,
		Memory:  64 * 1024, // 64 MiB
		Threads: 4,
		KeyLen:  32, // 32 bytes for AES-256
	}
}

// DeriveKey derives an encryption key from the master password using Argon2id.
func DeriveKey(masterPassword string, salt []byte, params KeyDerivationParams) []byte {
	return argon2.IDKey(
		[]byte(masterPassword),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLen,
	)
}

// Encrypt encrypts the plaintext using AES-256-GCM with the provided key.
// Returns the nonce, ciphertext, and an error if any.
func Encrypt(plaintext []byte, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return nonce, ciphertext, nil
}

// Decrypt decrypts the ciphertext using AES-256-GCM with the provided key and nonce.
// Returns the plaintext and an error if any.
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aesGCM.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:aesGCM.NonceSize()]
	ciphertext = ciphertext[aesGCM.NonceSize():]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptToString encrypts the plaintext and returns the result as a base64-encoded string.
func EncryptToString(plaintext []byte, key []byte) (string, error) {
	_, ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptFromString decrypts a base64-encoded ciphertext string.
func DecryptFromString(ciphertextStr string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return nil, err
	}
	return Decrypt(ciphertext, key)
}
