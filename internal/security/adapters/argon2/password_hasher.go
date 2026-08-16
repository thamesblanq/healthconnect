package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/thamesblanq/healthconnect/internal/security/ports"
)

const (
	time    uint32 = 3
	memory  uint32 = 64 * 1024
	threads uint8  = 4
	keyLen  uint32 = 32
	saltLen uint32 = 16
)

type PasswordHasher struct{}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	salt := make([]byte, saltLen)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		keyLen,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory,
		time,
		threads,
		encodedSalt,
		encodedHash,
	)

	return encoded, nil
}

func (h *PasswordHasher) Compare(
	storedHash string,
	password string,
) error {
	if storedHash == "" {
		return errors.New("stored password hash cannot be empty")
	}

	if password == "" {
		return errors.New("password cannot be empty")
	}

	parts := strings.Split(storedHash, "$")

	if len(parts) != 6 {
		return errors.New("invalid password hash format")
	}

	if parts[1] != "argon2id" {
		return errors.New("unsupported password hashing algorithm")
	}

	var parsedMemory uint32
	var parsedTime uint32
	var parsedThreads uint8

	_, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&parsedMemory,
		&parsedTime,
		&parsedThreads,
	)

	if err != nil {
		return errors.New("invalid Argon2 parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return errors.New("invalid password salt")
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return errors.New("invalid password hash")
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		parsedTime,
		parsedMemory,
		parsedThreads,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(actualHash, expectedHash) != 1 {
		return errors.New("invalid password")
	}

	return nil
}

// Compile-time check that PasswordHasher implements ports.PasswordHasher.
var _ ports.PasswordHasher = (*PasswordHasher)(nil)
