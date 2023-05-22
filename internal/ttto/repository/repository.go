package repository

import (
	"database/sql"
	"errors"

	"github.com/qerdcv/ttto/internal/ttto/repository/storage"
)

var (
	ErrUniqueViolation = errors.New("unique violation")
	ErrNotFound        = errors.New("not found")
)

type Repository struct {
	q *storage.Queries
}

func New(db *sql.DB) *Repository {
	return &Repository{
		q: storage.New(db),
	}
}
