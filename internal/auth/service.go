package auth

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/lckrugel/file-server/internal/users"
)

type UserProvider interface {
	Create(ctx context.Context, data *users.UserStore) (*users.User, error)
	FindByEmail(ctx context.Context, email string) (*users.User, error)
}

type AuthService struct {
	credRepo CredentialRepository
	userProv UserProvider
	hasher   *hasher
}

func NewAuthService(cRepo CredentialRepository, uProvider UserProvider) *AuthService {
	return &AuthService{
		credRepo: cRepo,
		userProv: uProvider,
		hasher:   newHasher(),
	}
}

func (s *AuthService) Register(
	ctx context.Context, email, name, password string,
) (string, error) {
	user, err := s.userProv.Create(ctx, &users.UserStore{
		Email: email,
		Name:  name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to store user: %w", err)
	}

	if err := validatePassword(password); err != nil {
		return "", err
	}

	passwordHash, err := s.hasher.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	cred := NewCredential(user.ID, passwordHash)
	if err := s.credRepo.Insert(ctx, cred); err != nil {
		return "", fmt.Errorf("failed to store user credentials: %w", err)
	}

	token, err := GenerateJWT(cred.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}
	return token, nil
}

func (s *AuthService) Login(
	ctx context.Context, email, password string,
) (*users.User, string, error) {
	user, err := s.userProv.FindByEmail(ctx, email)
	if errors.Is(err, users.ErrUserNotFound) {
		return nil, "", ErrInvalidCredentials
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to find user")
	}

	cred, err := s.credRepo.GetByUserID(ctx, user.ID)
	if errors.Is(err, ErrCredentialsNotFound) {
		return nil, "", ErrInvalidCredentials
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to find credentials")
	}

	match, err := s.hasher.VerifyPassword(password, cred.PasswordHash)
	if err != nil {
		return nil, "", fmt.Errorf("failed to verify password: %w", err)
	}
	if !match {
		return nil, "", ErrInvalidCredentials
	}

	token, err := GenerateJWT(cred.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT: %w", err)
	}
	return user, token, nil
}

func validatePassword(plainPassword string) error {
	if len(plainPassword) < 6 || len(plainPassword) > 64 {
		return ErrPasswordSize
	}
	if !utf8.ValidString(plainPassword) {
		return ErrPasswordInvalid
	}
	return nil
}
