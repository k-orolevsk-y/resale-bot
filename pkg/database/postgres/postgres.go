package postgres

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/config"
	"github.com/k-orolevsk-y/resale-bot/pkg/migrations"
)

type PgSQL interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	io.Closer

	Beginx() (TxPgSQL, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContextWithReturnID(ctx context.Context, query string, args ...interface{}) (interface{}, error)
}

type TxPgSQL interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	driver.Tx

	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContextWithReturnID(ctx context.Context, query string, args ...interface{}) (interface{}, error)
}

type postgresDatabase struct {
	*sqlx.DB
}

type postgresTx struct {
	*sqlx.Tx
}

func New() (PgSQL, error) {
	db, err := sqlx.Open("pgx/v5", config.Config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	if config.Config.MigrationsFlag {
		if err = migrations.ApplyMigrations(db); err != nil {
			return nil, fmt.Errorf("migrations.ApplyMigrations: %w", err)
		}
	}

	return &postgresDatabase{db}, err
}

func (db *postgresDatabase) Beginx() (TxPgSQL, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	return &postgresTx{
		Tx: tx,
	}, nil
}

func (db *postgresDatabase) ExecContextWithReturnID(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	query = fmt.Sprintf("%s RETURNING id", query)

	var id interface{}
	row := db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&id)
	return id, err
}

func (tx *postgresTx) ExecContextWithReturnID(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	query = fmt.Sprintf("%s RETURNING id", query)

	var id interface{}
	row := tx.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&id)
	return id, err
}
