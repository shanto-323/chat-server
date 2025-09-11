package logs

import (
	"log/slog"
	"os"
	"time"
)

const (
	service        = "client-service"
	env            = "dev"
	metricsTypeApi = "api_metrics"
	metricsTypeDb  = "db_metrics"
)

type InfoMetrics struct {
	MetricsType string
	Endpint     *string
	Status      int
	Latency     time.Duration
}

func NewInfoMetrics(endpoint *string, status int, latency time.Duration) *InfoMetrics {
	infoMetrics := InfoMetrics{}
	infoMetrics.MetricsType = metricsTypeDb
	if endpoint != nil {
		infoMetrics.MetricsType = metricsTypeApi
		infoMetrics.Endpint = endpoint
	}

	infoMetrics.Status = status
	infoMetrics.Latency = latency

	return &infoMetrics
}

type Logger interface {
	Info(value string)
	InfoMetrics(value *InfoMetrics)
	Error(err error)
	Warn(value any)
}

type slogLogger struct {
	logger *slog.Logger
}

func NewLogger(origin string) Logger {
	f, _ := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)

	logger := slog.New(slog.NewJSONHandler(f, nil))
	logger = logger.With(
		slog.Group("context",
			slog.String("service", service),
			slog.String("env", env),
			slog.String("origin", origin),
		),
	)
	return &slogLogger{
		logger: logger,
	}
}

func (l *slogLogger) Info(value string) {
	l.logger.Info("nornal", "msg", value)
}

func (l *slogLogger) InfoMetrics(value *InfoMetrics) {
	attrs := []any{
		slog.Int("status", value.Status),
		slog.Int64("latency_ms", value.Latency.Milliseconds()),
	}

	if value.Endpint != nil {
		attrs = append(attrs, slog.String("endpoint", *value.Endpint))
	}

	l.logger.Info(value.MetricsType, attrs...)
}

func (l *slogLogger) Error(err error) {
	l.logger.Error("normal", "error", err)
}

func (l *slogLogger) Warn(value any) {
	l.logger.Warn("normal", "warn", value)
}
