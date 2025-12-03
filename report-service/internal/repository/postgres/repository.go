package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Repository struct {
	cfg    *Config
	db     *pgxpool.Pool
	logger Logger
}

func NewRepository(cfg *Config, logger Logger) *Repository {
	return &Repository{
		cfg:    cfg,
		logger: logger,
	}
}

func (r *Repository) Run(_ context.Context) {}

func (r *Repository) Stop() {
	r.db.Close()
}

func (r *Repository) Init() error {
	var err error
	r.db, err = newPostgresDB(r.cfg)
	if err != nil {
		return err
	}

	return nil
}
