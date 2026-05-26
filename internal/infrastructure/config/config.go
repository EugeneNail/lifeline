package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App      App
	Database Database
}

type App struct {
	Name string
	Port int
}

type Database struct {
	Postgres Postgres
}

type Postgres struct {
	Name     string
	Port     int
	Host     string
	User     string
	Password string
}

// Load builds Config from environment variables and returns an error when a required value is missing or invalid.
func Load() (Config, error) {
	appName, err := readString("APP_NAME")
	if err != nil {
		return Config{}, fmt.Errorf("loading config.app.name: %w", err)
	}

	appPort, err := readInt("APP_PORT")
	if err != nil {
		return Config{}, fmt.Errorf("loading config.app.port: %w", err)
	}

	postgresName, err := readString("DATABASE_POSTGRES_NAME")
	if err != nil {
		return Config{}, fmt.Errorf("loading config.database.postgres.name: %w", err)
	}

	postgresPort, err := readInt("DATABASE_POSTGRES_PORT")
	if err != nil {
		return Config{}, fmt.Errorf("loading config.database.postgres.port: %w", err)
	}

	postgresHost, err := readString("DATABASE_POSTGRES_HOST")
	if err != nil {
		return Config{}, fmt.Errorf("loading config.database.postgres.host: %w", err)
	}

	postgresUser, err := readString("DATABASE_POSTGRES_USER")
	if err != nil {
		return Config{}, fmt.Errorf("loading config.database.postgres.user: %w", err)
	}

	postgresPassword, err := readString("DATABASE_POSTGRES_PASSWORD")
	if err != nil {
		return Config{}, fmt.Errorf("loading config.database.postgres.password: %w", err)
	}

	return Config{
		App: App{
			Name: appName,
			Port: appPort,
		},
		Database: Database{
			Postgres: Postgres{
				Name:     postgresName,
				Port:     postgresPort,
				Host:     postgresHost,
				User:     postgresUser,
				Password: postgresPassword,
			},
		},
	}, nil
}

// readString returns a non-empty environment variable value or an error when the variable is missing or empty.
func readString(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("env %s is not set", key)
	}

	if value == "" {
		return "", fmt.Errorf("env %s is empty", key)
	}

	return value, nil
}

// readInt parses an environment variable into int and returns an error when the variable is missing, empty, or not a valid integer.
func readInt(key string) (int, error) {
	value, err := readString(key)
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parsing env %s value %q into integer: %w", key, value, err)
	}

	return number, nil
}
