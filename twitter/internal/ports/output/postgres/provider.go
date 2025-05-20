package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/twitter-tq/vinofsteel/pkg/database"
)

type PostgresProvider struct {
	connStr string
	db      *sql.DB
}

func NewPostgresDatabaseProvider() database.SQLProvider {
	return newPostgresProvider(
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)
}


func (p *PostgresProvider) GetConnection(ctx context.Context) (*sql.DB, error) {
	if p.db != nil {
		if err := p.db.PingContext(ctx); err != nil {
			err := p.db.Close()
			return nil, err
		}

		return p.db, nil
	}

	db, err := sql.Open("postgres", p.connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening DB connection: %w", err)
	}

	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		if err := db.Close(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	p.db = db
	return db, nil
}

func (p *PostgresProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}

	return nil
}

// Utilities
func newPostgresProvider(user, password, host, port, dbName string) *PostgresProvider {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	return &PostgresProvider{
		connStr: connStr,
	}
}
