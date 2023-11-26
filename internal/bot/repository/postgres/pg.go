package postgres

import "github.com/k-orolevsk-y/resale-bot/pkg/database/postgres"

type Pg struct {
	db postgres.PgSQL
}

func New(db postgres.PgSQL) *Pg {
	return &Pg{db: db}
}
