package service

import (
	"github.com/qerdcv/ttto/internal/auth"
	"github.com/qerdcv/ttto/internal/conf"
	"github.com/qerdcv/ttto/internal/ttto/repository"
)

type Service struct {
	repo      *repository.Repository
	tokenizer *auth.JWTTokenizer

	cfg conf.App
}

func New(repo *repository.Repository, tokenizer *auth.JWTTokenizer) *Service {
	return &Service{
		repo:      repo,
		tokenizer: tokenizer,
	}
}
