package postgres

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/EugeneNail/lifeline/internal/infrastructure/config"
	_ "github.com/lib/pq"
)

// Connect builds a PostgreSQL connection string from config, opens the database, and returns a ready-to-use SQL handle.
func Connect(configuration config.Postgres) (*sql.DB, error) {
	connectionString := (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(configuration.User, configuration.Password),
		Host:     fmt.Sprintf("%s:%d", configuration.Host, configuration.Port),
		Path:     configuration.Name,
		RawQuery: "sslmode=disable",
	}).String()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("opening postgres database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("pinging postgres database: %w", err)
	}

	return db, nil
}
