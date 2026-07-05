package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	UserID       uuid.UUID
	PasswordHash string
	UpdatedAt    time.Time
	LastUsedAt   time.Time
}

func NewCredential(userID uuid.UUID, pass string) *Credential {
	return &Credential{
		UserID:       userID,
		PasswordHash: pass,
		UpdatedAt:    time.Now(),
	}
}

type CredentialRepository interface {
	Insert(ctx context.Context, cred *Credential) error
	Update(ctx context.Context, cred *Credential) (*Credential, error)
	GetByUserID(ctx context.Context, id uuid.UUID) (*Credential, error)
}

var (
	ErrCredentialsNotFound = errors.New("user credentials not found")
	ErrPasswordSize        = errors.New("invalid password size")
	ErrPasswordInvalid     = errors.New("password contains invalid characters")
	ErrInvalidCredentials  = errors.New("email or password doesn't match")
)
