package delivery

import (
	"context"
	"net/http"
	"tasks-service/internal/application"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Config struct {
	Port         string        `env:"PORT"`
	ReadTimeOut  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeOut time.Duration `env:"WRITE_TIMEOUT"`
}

type Handler struct {
	cfg        *Config
	services   *application.Service
	httpServer *http.Server
	router     *gin.Engine
	logger     Logger
}

func NewHandler(cfg *Config, service *application.Service, logger Logger) *Handler {
	return &Handler{
		cfg:      cfg,
		services: service,
		logger:   logger,
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

		task := api.Group("/tasks", h.CheckJWT)
		{
			task.POST("/", h.createTask)
			task.GET("/:id", h.getTasks)
			task.PUT("/:id", h.updateTask)
			task.DELETE("/:id", h.deleteTask)
		}

	}

	h.router = router
	return nil
}
