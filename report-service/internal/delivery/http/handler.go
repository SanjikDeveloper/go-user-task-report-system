package http

import (
	"context"
	"html/template"
	"net/http"
	"report-service/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

// Мы держим интерфейс в месте его использования
type Service interface {
	ReportByRange(ctx context.Context, userID int64, username string, start, end time.Time) (models.ReportVM, error)
	PrepareReportData(username string, start, end time.Time, tasks []models.Task) models.ReportVM
}

type Config struct {
	Port         string        `env:"PORT"`
	ReadTimeOut  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeOut time.Duration `env:"WRITE_TIMEOUT"`
}

type Handler struct {
	cfg        *Config
	services   Service
	httpServer *http.Server
	router     *gin.Engine
	logger     Logger
}

func NewHandler(cfg *Config, service Service, logger Logger) *Handler {
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
	// TODO: добавить api/v1

	api := router.Group("/api/v1")
	{

		report := api.Group("/report")
		{
			report.GET("/", h.GetReport)
		}

	}

	tmpl, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		return err
	}

	router.SetHTMLTemplate(tmpl)

	h.router = router
	return nil
}
