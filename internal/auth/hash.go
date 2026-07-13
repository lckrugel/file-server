package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type hasher struct {
	memory      uint32 // In KiB
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func newHasher() *hasher {
	return &hasher{
		memory:      19 * 1024,
		iterations:  2,
		parallelism: 1,
		saltLength:  16,
		keyLength:   32,
	}
}

func (h *hasher) hashPassword(plainPassword string) (string, error) {
	salt := make([]byte, h.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(plainPassword),
		salt,
		h.iterations,
		h.memory,
		h.parallelism,
		h.keyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$%s$%s", b64Salt, b64Hash)

	return encoded, nil
}

func (h *hasher) verifyPassword(givenPassword, encodedHash string) (bool, error) {
	salt, hash, err := h.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	givenHash := argon2.IDKey(
		[]byte(givenPassword),
		salt,
		h.iterations,
		h.memory,
		h.parallelism,
		h.keyLength,
	)

	return subtle.ConstantTimeCompare(hash, givenHash) == 1, nil
}

func (h *hasher) decodeHash(encodedHash string) ([]byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) < 3 {
		return nil, nil, errors.New("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, err
	}

	return salt, hash, nil
}
