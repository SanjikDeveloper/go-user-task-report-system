package application

import (
	"context"
	"user-service/internal/models"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, username, password string) (models.User, error)
}

type Service struct {
	repo   Repository
	logger Logger
	auth   *AuthConfig
}

type AuthConfig struct {
	JWTSigningKey string
	PasswordSalt  string
}

func NewService(repo Repository, logger Logger, auth *AuthConfig) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
		auth:   auth,
	}
}

func (s *Service) Init() error {
	return nil
}

func (s *Service) Run(_ context.Context) {
}

func (s *Service) Stop() {
}
