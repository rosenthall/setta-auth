package app

import (
	"auth_service/internal/configuration"
	"auth_service/internal/repository"
	"auth_service/internal/server"
	service "auth_service/internal/services/auth"
	"context"
	"log"
)

func Run(config *configuration.AuthServiceConfig) {

	logger, loggerError := getLogger(config)
	if loggerError != nil {
		log.Fatalf("Error while getting setting up the logger : " + loggerError.Error())
	}

	logger.Info("Logger has successfully initialized")

	logger.Info("Getting signing keys..")

	// Getting RSA keys
	publicKey, pubKeyError := getPublicKeyFromFile(config.PublicKeyPath)
	if pubKeyError != nil {
		panic(pubKeyError)
	}
	logger.Debugf("Public key size : %#v", publicKey.Size())

	privateKey, privateKeyError := getPrivateKeyFromFile(config.PrivateKeyPath)
	if privateKeyError != nil {
		panic(privateKeyError)
	}
	logger.Debugf("Private key size : %#v ", privateKey.Size())

	// Connecting to the redis server
	logger.Info("Connecting to redis server..")
	redisClient, redisErr := getRedisConnection(context.Background(), config.RedisServerIp, config.RedisPassword, logger)
	if redisErr != nil {
		panic(redisErr)
	}
	defer redisClient.Close() // Closing connection there

	logger.Info("Successfully connected!")

	// Initializing our repository
	redisRepository := repository.NewRedisRepository(redisClient, *logger)

	jwtAuthService := service.NewJWTAuthService(config, logger, *redisRepository, privateKey, publicKey)

	grpcServer := server.NewServer(config, logger, jwtAuthService)
	err := grpcServer.Run()
	if err != nil {
		panic(err)
	}
}
