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
	cfg          *Config
	db           *pgxpool.Pool
	TaskPostgres *TaskPostgres
	logger       Logger
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
	// Если уже инициализирован, пропускаем
	if r.db != nil && r.TaskPostgres != nil {
		return nil
	}

	var err error
	r.db, err = newPostgresDB(r.cfg)
	if err != nil {
		return err
	}
	r.TaskPostgres = NewTaskPostgres(r.db)

	// Автоматически создаем таблицы, если их нет
	if err := r.migrate(); err != nil {
		r.logger.Error("failed to migrate database: %v", err)
		return err
	}

	return nil
}

func (r *Repository) migrate() error {
	ctx := context.Background()

	r.logger.Info("Starting database migration...")
	r.logger.Info("Database name: %s", r.cfg.DBName)

	// Проверяем существование таблицы users
	var usersExists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users')",
	).Scan(&usersExists)
	if err != nil {
		r.logger.Error("Failed to check users table: %v", err)
		return err
	}

	if !usersExists {
		r.logger.Info("Users table not found, creating...")
		_, err = r.db.Exec(ctx, `
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				username VARCHAR(255) NOT NULL UNIQUE,
				password VARCHAR(255) NOT NULL,
				role VARCHAR(100) NOT NULL DEFAULT 'user',
				date_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
		`)
		if err != nil {
			r.logger.Error("Failed to create users table: %v", err)
			return err
		}
		r.logger.Info("Users table created successfully")
	} else {
		r.logger.Info("Users table already exists")
	}

	// Проверяем существование таблицы tasks
	var tasksExists bool
	err = r.db.QueryRow(ctx,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'tasks')",
	).Scan(&tasksExists)
	if err != nil {
		r.logger.Error("Failed to check tasks table: %v", err)
		return err
	}

	if !tasksExists {
		r.logger.Info("Tasks table not found, creating...")
		_, err = r.db.Exec(ctx, `
			CREATE TABLE tasks (
				id SERIAL PRIMARY KEY,
				user_id INTEGER NOT NULL,
				title VARCHAR(255) NOT NULL,
				description TEXT NOT NULL,
				status VARCHAR(50) NOT NULL DEFAULT 'new',
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			);
			CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
			CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
			CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
		`)
		if err != nil {
			r.logger.Error("Failed to create tasks table: %v", err)
			return err
		}
		r.logger.Info("Tasks table created successfully")
	} else {
		r.logger.Info("Tasks table already exists")
	}

	r.logger.Info("Database migration completed successfully")
	return nil
}
