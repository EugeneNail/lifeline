package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/EugeneNail/lifeline/internal/infrastructure/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

const migrationsDir = "internal/infrastructure/postgres/migrations"

var migrationNameSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

// main runs the migrator command and exits with a non-zero status when execution fails.
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "migrator error: %v\n", err)
		os.Exit(1)
	}
}

// run dispatches the requested migrator command and returns an error when the command is unknown or fails.
func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required: create, migrate, rollback")
	}

	switch args[0] {
	case "create":
		return runCreate(args[1:])
	case "migrate":
		return runMigrate()
	case "rollback":
		return runRollback()
	default:
		return fmt.Errorf("unsupported command %q", args[0])
	}
}

// runCreate creates paired up and down migration files and returns an error when the migration name is invalid or file creation fails.
func runCreate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("create command requires migration name")
	}

	migrationName, err := normalizeMigrationName(strings.Join(args, " "))
	if err != nil {
		return fmt.Errorf("creating migration files: %w", err)
	}

	timestamp := time.Now().UTC().Format("20060102150405")
	filePrefix := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s", timestamp, migrationName))

	if err := os.MkdirAll(migrationsDir, 0o755); err != nil {
		return fmt.Errorf("creating migrations directory %q: %w", migrationsDir, err)
	}

	if err := createEmptyFile(filePrefix + ".up.sql"); err != nil {
		return fmt.Errorf("creating up migration file: %w", err)
	}

	if err := createEmptyFile(filePrefix + ".down.sql"); err != nil {
		return fmt.Errorf("creating down migration file: %w", err)
	}

	return nil
}

// runMigrate ensures the target database exists, applies all pending migrations, and returns an error when the operation fails.
func runMigrate() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading migrator config: %w", err)
	}

	if err := ensureDatabaseReady(cfg); err != nil {
		return fmt.Errorf("preparing database for migrate command: %w", err)
	}

	migrator, err := migrate.New(buildMigrationsURL(), buildDatabaseURL(cfg, cfg.Database.Postgres.Name))
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	defer closeMigrator(migrator)

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}

// runRollback ensures the target database exists, rolls back the last migration, and returns an error when the operation fails.
func runRollback() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading migrator config: %w", err)
	}

	if err := ensureDatabaseReady(cfg); err != nil {
		return fmt.Errorf("preparing database for rollback command: %w", err)
	}

	migrator, err := migrate.New(buildMigrationsURL(), buildDatabaseURL(cfg, cfg.Database.Postgres.Name))
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	defer closeMigrator(migrator)

	if err := migrator.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("rolling back last migration: %w", err)
	}

	return nil
}

// ensureDatabaseReady verifies server connectivity, creates the target database when it does not exist, and returns an error when preparation fails.
func ensureDatabaseReady(cfg config.Config) error {
	adminDatabase, err := openDatabase(cfg, "postgres")
	if err != nil {
		return fmt.Errorf("opening postgres maintenance database: %w", err)
	}

	defer closeDatabase(adminDatabase, "postgres maintenance")

	if err := pingDatabase(adminDatabase, "postgres maintenance"); err != nil {
		return fmt.Errorf("checking postgres maintenance connectivity: %w", err)
	}

	if err := ensureDatabaseExists(adminDatabase, cfg.Database.Postgres.Name); err != nil {
		return fmt.Errorf("ensuring postgres database %q exists: %w", cfg.Database.Postgres.Name, err)
	}

	targetDatabase, err := openDatabase(cfg, cfg.Database.Postgres.Name)
	if err != nil {
		return fmt.Errorf("opening postgres database %q: %w", cfg.Database.Postgres.Name, err)
	}

	defer closeDatabase(targetDatabase, cfg.Database.Postgres.Name)

	if err := pingDatabase(targetDatabase, cfg.Database.Postgres.Name); err != nil {
		return fmt.Errorf("checking postgres database %q connectivity: %w", cfg.Database.Postgres.Name, err)
	}

	return nil
}

// ensureDatabaseExists creates the target database when it is absent and returns an error when the existence check or creation fails.
func ensureDatabaseExists(db *sql.DB, databaseName string) error {
	exists, err := databaseExists(db, databaseName)
	if err != nil {
		return fmt.Errorf("checking database %q existence: %w", databaseName, err)
	}

	if exists {
		return nil
	}

	query := fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(databaseName))
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("creating database %q: %w", databaseName, err)
	}

	return nil
}

// databaseExists reports whether the target database exists and returns an error when the lookup fails.
func databaseExists(db *sql.DB, databaseName string) (bool, error) {
	var exists bool

	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", databaseName).Scan(&exists); err != nil {
		return false, fmt.Errorf("querying database %q existence: %w", databaseName, err)
	}

	return exists, nil
}

// openDatabase opens a SQL connection for the provided database name and returns an error when the driver cannot initialize it.
func openDatabase(cfg config.Config, databaseName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", buildConnectionString(cfg, databaseName))
	if err != nil {
		return nil, fmt.Errorf("opening postgres connection for database %q: %w", databaseName, err)
	}

	return db, nil
}

// pingDatabase verifies that the provided database connection is reachable and returns an error when the ping fails.
func pingDatabase(db *sql.DB, databaseName string) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("pinging postgres database %q: %w", databaseName, err)
	}

	return nil
}

// closeDatabase closes the SQL connection and writes a warning to stderr when closing fails.
func closeDatabase(db *sql.DB, databaseName string) {
	if err := db.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "migrator warning: closing database %q: %v\n", databaseName, err)
	}
}

// closeMigrator closes the migrate instance and writes a warning to stderr when closing fails.
func closeMigrator(migrator *migrate.Migrate) {
	sourceErr, databaseErr := migrator.Close()
	if sourceErr != nil {
		fmt.Fprintf(os.Stderr, "migrator warning: closing migration source: %v\n", sourceErr)
	}

	if databaseErr != nil {
		fmt.Fprintf(os.Stderr, "migrator warning: closing migration database: %v\n", databaseErr)
	}
}

// buildMigrationsURL returns the file URL pointing to the migrations directory used by golang-migrate.
func buildMigrationsURL() string {
	return fmt.Sprintf("file://%s", migrationsDir)
}

// buildDatabaseURL returns a PostgreSQL URL suitable for golang-migrate for the provided database name.
func buildDatabaseURL(cfg config.Config, databaseName string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		databaseName,
	)
}

// buildConnectionString returns a PostgreSQL DSN suitable for database/sql for the provided database name.
func buildConnectionString(cfg config.Config, databaseName string) string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		databaseName,
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
	)
}

// normalizeMigrationName converts a free-form migration title into a snake_case file name and returns an error when no valid characters remain.
func normalizeMigrationName(name string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = migrationNameSanitizer.ReplaceAllString(normalized, "_")
	normalized = strings.Trim(normalized, "_")

	if normalized == "" {
		return "", fmt.Errorf("migration name %q does not contain letters or digits", name)
	}

	return normalized, nil
}

// createEmptyFile creates an empty file and returns an error when the file already exists or cannot be created.
func createEmptyFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("opening file %q for creation: %w", path, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("closing newly created file %q: %w", path, err)
	}

	return nil
}
