package users

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
}

func NewUser(email, name string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

type UserRepository interface {
	Insert(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, u *User) (*User, error)
	List(ctx context.Context) ([]*User, error)
}

// Domain errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailInUse   = errors.New("email already exists")
)

type UserStore struct {
	Email string
	Name  string
}

type UserUpdate struct {
	Name string
}
