package main

import (
	"context"
	"log/slog"
	"report-service/internal/application"
	delivery "report-service/internal/delivery/http"
	"report-service/internal/repository/postgres"
	"report-service/pkg/config"
	"report-service/pkg/logger"
	service "report-service/pkg/services"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Repo   postgres.Config `envPrefix:"REPO_"`
	Logger logger.Config   `envPrefix:"LOGGER_"`
	Http   delivery.Config `envPrefix:"HTTP_"`
	Auth   config.AuthConfig
}

// TODO: использовать линтеры, прогуглить про golangcilint конфиг и записать в свои проекты конфиги
func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("error loading env variables: %s", err.Error())
		return
	}

	cfg := Config{}
	if err := config.ReadEnvConfig(&cfg); err != nil {
		slog.Error("error initializing configs: %s", err.Error())
		return
	}

	log := logger.NewLogger(&cfg.Logger)

	repos := postgres.NewRepository(&cfg.Repo, log)
	authCfg := &application.AuthConfig{
		JWTSigningKey: cfg.Auth.JWTSigningKey,
	}
	services := application.NewService(repos, log, authCfg)
	handlers := delivery.NewHandler(&cfg.Http, services, log)

	srv := service.NewManager(log)
	srv.AddService(
		repos,
		handlers,
	)
	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		err := errors.Wrap(err, "srv.Run(...) err:")
		log.Error(err.Error())
		return
	}

	log.Info("User-service Started")
}
