package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/urfave/cli/v2"

	"github.com/qerdcv/ttto/internal/conf"
	"github.com/qerdcv/ttto/internal/ttto/repository/migrations"
)

func runMigration(c *cli.Context) error {
	cfg, err := conf.New()
	if err != nil {
		return fmt.Errorf("conf new: %w", err)
	}

	d, err := iofs.New(migrations.Migrations, ".")
	if err != nil {
		return fmt.Errorf("iofs new: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, cfg.DB.DSN())
	if err != nil {
		return fmt.Errorf("migrate new with source instance: %w", err)
	}

	v, isDirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("migration version: %w", err)
	}

	log.Printf("before migration current version %d, is dirty %v", v, isDirty)

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up: %w", err)
	}

	v, isDirty, err = m.Version()
	log.Printf("after migration current version %d, is dirty %v", v, isDirty)

	return nil
}
