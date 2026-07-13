package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, data *UserStore) (*User, error) {
	existing, err := s.repo.GetByEmail(ctx, data.Email)
	if err != nil && err != ErrUserNotFound {
		return nil, fmt.Errorf("failed to fetch existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailInUse
	}
	user := NewUser(data.Email, data.Name)
	if err := s.repo.Insert(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to insert new user: %w", err)
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, id string, data *UserUpdate) (*User, error) {
	user, err := s.FindById(ctx, id)
	if errors.Is(err, ErrUserNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	user.Name = data.Name
	user, err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (s *UserService) FindById(ctx context.Context, id string) (*User, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}
	user, err := s.repo.GetById(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
