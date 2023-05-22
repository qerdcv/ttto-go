package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/ttto/repository/storage"
)

func (r *Repository) CreateUser(ctx context.Context, user domain.User) error {
	if _, err := r.q.CreateUser(ctx, storage.CreateUserParams{
		Username: user.Username,
		Password: user.Password,
	}); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return ErrUniqueViolation
		}

		return fmt.Errorf("queries create user: %w", err)
	}

	return nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	u, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, ErrNotFound
		}

		return domain.User{}, fmt.Errorf("queries get user by username: %w", err)
	}

	return domain.User(u), nil
}
