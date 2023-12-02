package postgres

import (
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/k-orolevsk-y/resale-bot/pkg/database/postgres"
)

type Pg struct {
	db postgres.PgSQL
}

type Begin interface {
	Begin() (postgres.PgSQL, error)
}

func New(db postgres.PgSQL) *Pg {
	return &Pg{db: db}
}

func (pg *Pg) Close() error {
	closedDB, ok := pg.db.(io.Closer)
	if !ok {
		return fmt.Errorf("error this is not io.Closer")
	}

	return closedDB.Close()
}

func (pg *Pg) Commit() error {
	tx, ok := pg.db.(driver.Tx)
	if !ok {
		return fmt.Errorf("error this is not driver.Tx")
	}

	return tx.Commit()
}

func (pg *Pg) Rollback() error {
	tx, ok := pg.db.(driver.Tx)
	if !ok {
		return fmt.Errorf("error this is not driver.Tx")
	}

	return tx.Rollback()
}

func (pg *Pg) WithTx() (*Pg, error) {
	db, ok := pg.db.(Begin)
	if !ok {
		return nil, fmt.Errorf("error this is not driver.Conn")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return &Pg{db: tx}, nil
}
