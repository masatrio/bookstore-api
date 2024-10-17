package postgresql

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq" // Import PostgreSQL driver
	"github.com/masatrio/bookstore-api/config"
)

// NewDatabase initializes a new database connection and returns it.
func NewDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConnection)
	db.SetMaxOpenConns(cfg.MaxActiveConnection)
	db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)

	return db, nil
}
