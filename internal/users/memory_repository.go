package users

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	users map[uuid.UUID]*User
	mu    sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[uuid.UUID]*User),
	}
}

// Validates that MemoryRepository implementes UserRepository
var _ UserRepository = (*MemoryRepository)(nil)

func (memo *MemoryRepository) Insert(ctx context.Context, u *User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	memo.users[u.ID] = u
	return nil
}

func (memo *MemoryRepository) GetById(ctx context.Context, id uuid.UUID) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	for _, user := range memo.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (memo *MemoryRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	for _, user := range memo.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (memo *MemoryRepository) List(ctx context.Context) ([]*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	var users []*User
	for _, file := range memo.users {
		users = append(users, file)

	}
	return users, nil
}

func (memo *MemoryRepository) Update(ctx context.Context, u *User) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	memo.users[u.ID] = u
	return u, nil
}
