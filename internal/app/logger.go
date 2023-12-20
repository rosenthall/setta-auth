package app

import (
	"auth_service/internal/configuration"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

func getLogger(cfg *configuration.AuthServiceConfig) (*zap.SugaredLogger, error) {
	var logLevel zapcore.Level

	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.InfoLevel // Default level is Info
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      false,
		Encoding:         "console", // or "json"
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()
	return sugar, nil
}
