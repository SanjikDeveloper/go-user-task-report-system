package application

import (
	"context"
	"report-service/internal/models"
	"time"
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

// TODO: интерфейсы хранить в месте использования
// все методы из repository нужны в интерфейсе внутри application
// все методы из application нужны в интерфейсе внутри delivery

type Repository interface {
	ListByRange(ctx context.Context, userId int64, start, end time.Time) ([]models.Task, error)
	ListByDate(ctx context.Context, userId int64, date time.Time) ([]models.Task, error)
}

type Service struct {
	repo   Repository
	logger Logger
	auth   *AuthConfig
}

func NewService(repo Repository, logger Logger, auth *AuthConfig) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
		auth:   auth,
	}
}

func (s *Service) GetJWTSigningKey() string {
	return s.auth.JWTSigningKey
}
