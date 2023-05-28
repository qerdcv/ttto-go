package service

import (
	"github.com/qerdcv/ttto/internal/auth"
	"github.com/qerdcv/ttto/internal/conf"
	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/eventst"
	"github.com/qerdcv/ttto/internal/ttto/repository"
)

type Service struct {
	repo      *repository.Repository
	tokenizer *auth.JWTTokenizer
	es        *eventst.EventStream[*domain.Game]

	cfg conf.App
}

func New(repo *repository.Repository, tokenizer *auth.JWTTokenizer, es *eventst.EventStream[*domain.Game], cfg conf.App) *Service {
	return &Service{
		repo:      repo,
		tokenizer: tokenizer,
		es:        es,
		cfg:       cfg,
	}
}
