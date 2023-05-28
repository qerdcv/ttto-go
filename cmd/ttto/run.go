package main

import (
	"database/sql"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/qerdcv/ttto/internal/auth"
	"github.com/qerdcv/ttto/internal/conf"
	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/eventst"
	"github.com/qerdcv/ttto/internal/ttto/repository"
	"github.com/qerdcv/ttto/internal/ttto/runner"
	"github.com/qerdcv/ttto/internal/ttto/runner/http"
	"github.com/qerdcv/ttto/internal/ttto/runner/metrics"
	"github.com/qerdcv/ttto/internal/ttto/service"
)

func run(c *cli.Context) error {
	cfg, err := conf.New()
	if err != nil {
		return fmt.Errorf("conf new: %w", err)
	}

	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}

	es := eventst.NewEventStream[*domain.Game]()
	repo := repository.New(db)
	tokenizer := auth.NewJWTTokenizer(cfg.Auth)
	svc := service.New(repo, tokenizer, es, cfg.App)

	if err = runner.Run(
		http.New(svc, es, cfg.HTTP),
		metrics.New(cfg.Metrics),
	); err != nil {
		return fmt.Errorf("runner run: %w", err)
	}

	return nil
}
