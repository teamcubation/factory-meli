package database

import (
	"context"
	"database/sql"
)

// DatabaseProvider is the abstraction for all SQL database use in the application
type SQLProvider interface {
	GetConnection(ctx context.Context) (*sql.DB, error)
	Close() error
}
