package main

import (
	"database/sql"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/qerdcv/ttto/internal/auth"
	"github.com/qerdcv/ttto/internal/conf"
	"github.com/qerdcv/ttto/internal/ttto/repository"
	"github.com/qerdcv/ttto/internal/ttto/runner/http"
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

	repo := repository.New(db)
	tokenizer := auth.NewJWTTokenizer(cfg.Auth)
	svc := service.New(repo, tokenizer)
	httpR := http.New(svc, cfg.HTTP)

	if err := httpR.Run(); err != nil {
		return fmt.Errorf("http runner run: %w", err)
	}

	return nil
}
