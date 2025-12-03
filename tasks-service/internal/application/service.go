package application

import (
	"context"
	"tasks-service/internal/models"
	"tasks-service/internal/repository/postgres"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type AuthConfig struct {
	JWTSigningKey string
}

type TaskTodo interface {
	CreateTask(ctx context.Context, task models.Task) (int, error)
	GetTaskById(ctx context.Context, id, userID int) (models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, id, userID int) error
}

type Service struct {
	TaskTodo
	logger Logger
	auth   *AuthConfig
}

func NewService(postgres *postgres.Repository, logger Logger, auth *AuthConfig) *Service {
	return &Service{
		TaskTodo: NewTaskService(postgres.TaskPostgres),
		logger:   logger,
		auth:     auth,
	}
}

func (s *Service) GetJWTSigningKey() string {
	return s.auth.JWTSigningKey
}
