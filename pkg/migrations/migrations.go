package migrations

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var sqlMigrationFiles embed.FS

func ApplyMigrations(db *sqlx.DB) error {
	goose.SetBaseFS(sqlMigrationFiles)
	goose.SetSequential(true)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(db.DB, "."); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}
