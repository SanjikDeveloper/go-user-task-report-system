package delivery

import (
	"context"
	"net/http"
	"time"
	"user-service/internal/models"

	"github.com/gin-gonic/gin"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Service interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(acessToken string) (int, error)
}

type Config struct {
	Port         string        `env:"PORT"`
	ReadTimeOut  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeOut time.Duration `env:"WRITE_TIMEOUT"`
}

type Handler struct {
	cfg        *Config
	service    Service
	httpServer *http.Server
	router     *gin.Engine
	logger     Logger
}

func NewHandler(cfg *Config, service Service, logger Logger) *Handler {
	return &Handler{
		cfg:     cfg,
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Run(_ context.Context) {

	h.httpServer = &http.Server{
		Addr:         ":" + h.cfg.Port,
		Handler:      h.router,
		ReadTimeout:  h.cfg.ReadTimeOut,
		WriteTimeout: h.cfg.WriteTimeOut,
	}

	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil {
			h.logger.Error("HTTP server stopped unexpectedly", "error", err)
			return
		}
	}()
}

func (h *Handler) Stop() {
	err := h.httpServer.Shutdown(context.Background())
	if err != nil {
		h.logger.Error("HTTP server stopped unexpectedly", "error", err)
	}
}

func (h *Handler) Init() error {
	router := gin.New()

	api := router.Group("/api/v1")
	{
		auth := api.Group("/user")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)

		}
	}

	h.router = router
	return nil
}
