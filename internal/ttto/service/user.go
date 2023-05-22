package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/ttto/repository"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *Service) CreateUser(ctx context.Context, user domain.User) error {
	if err := user.Validate(); err != nil {
		return newErrValidation(err)
	}

	user.HashPassword(s.cfg.Secret)
	if err := s.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUniqueViolation) {
			return ErrUserAlreadyExists
		}

		return fmt.Errorf("repo create user: %w", err)
	}

	return nil
}

func (s *Service) AuthorizeUser(ctx context.Context, user domain.User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", newErrValidation(err)
	}

	dbUser, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrUserNotFound
		}

		return "", fmt.Errorf("repo get user by username: %w", err)
	}

	user.HashPassword(s.cfg.Secret)
	if dbUser.Password != user.Password {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokenizer.Encode(dbUser)
	if err != nil {
		return "", fmt.Errorf("tokenizer encode: %w", err)
	}

	return token, nil
}

func (s *Service) DecodeToken(token string) (domain.User, error) {
	return s.tokenizer.Decode(token)
}
