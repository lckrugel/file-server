package auth

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	credentials map[uuid.UUID]*Credential
	mu          sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		credentials: make(map[uuid.UUID]*Credential),
	}
}

// Validates that MemoryRepository implementes CredentialRepository
var _ CredentialRepository = (*MemoryRepository)(nil)

func (memo *MemoryRepository) Insert(ctx context.Context, cred *Credential) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	memo.credentials[cred.UserID] = cred
	return nil
}

func (memo *MemoryRepository) Update(ctx context.Context, cred *Credential) (*Credential, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	memo.credentials[cred.UserID] = cred
	return cred, nil
}

func (memo *MemoryRepository) GetByUserID(ctx context.Context, id uuid.UUID) (*Credential, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	for _, cred := range memo.credentials {
		if cred.UserID == id {
			return cred, nil
		}
	}
	return nil, ErrCredentialsNotFound
}
