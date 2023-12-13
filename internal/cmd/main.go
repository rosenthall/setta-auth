package main

import (
	"auth_service/internal/configuration"
	"auth_service/internal/server"
	service "auth_service/internal/services"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"strings"
)

func main() {
	config, configError := configuration.ReadConfig("../config.toml")
	if configError != nil {
		log.Fatalf("Error while getting configuration : " + configError.Error())
	}

	logger, loggerError := getLogger(config)
	if loggerError != nil {
		log.Fatalf("Error while getting setting up the logger : " + loggerError.Error())
	}

	logger.Info("Logger has successfully initialized")

	jwtAuthService, _ := service.NewJWTAuthService(config, logger)

	grpcServer := server.NewServer(config, logger, jwtAuthService)
	err := grpcServer.Run()
	if err != nil {
		panic(err)
	}
}

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
