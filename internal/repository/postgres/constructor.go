package postgres

import (
	"context"
	"fmt"

	"user-crud/internal/config"

	"github.com/bigmate/closer"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

//NewClient is the PostgresDB client
func NewClient(ctx context.Context, config *config.Config) (*Client, error) {
	dsn := fmt.Sprintf(
		config.Postgres.DSN,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.DBName,
		config.Postgres.Port)

	db, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	closer.Add(db.Close)

	return &Client{db: db}, nil
}
