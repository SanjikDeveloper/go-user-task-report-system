package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	DBName   string `env:"NAME"`
	SSLMode  string `env:"SSLMODE"`
}

func newPostgresDB(cfg *Config) (*pgxpool.Pool, error) {
	fmt.Printf("[DB] Connecting to database: %s@%s:%s/%s\n", cfg.Username, cfg.Host, cfg.Port, cfg.DBName)

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w. Check your database configuration (host=%s, port=%s, username=%s, dbname=%s)",
			err, cfg.Host, cfg.Port, cfg.Username, cfg.DBName)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database '%s': %w. Make sure database exists and is accessible", cfg.DBName, err)
	}

	fmt.Printf("[DB] Successfully connected to database: %s\n", cfg.DBName)
	return pool, nil
}
