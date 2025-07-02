package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerInterface interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Sync() error
}

type ZapLogger struct {
	logger *zap.Logger
}

var _ loggerInterface = (*ZapLogger)(nil)

func NewLogger() (*ZapLogger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger,
	}, nil
}

func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}
