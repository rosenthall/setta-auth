package app

import (
	"auth_service/internal/configuration"
	"auth_service/internal/server"
	service "auth_service/internal/services/auth"
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

	jwtAuthService, _ := service.NewJWTAuthService(config, logger, privateKey, publicKey)

	grpcServer := server.NewServer(config, logger, jwtAuthService)
	err := grpcServer.Run()
	if err != nil {
		panic(err)
	}
}
