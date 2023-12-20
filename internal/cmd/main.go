package main

import (
	"auth_service/internal/app"
	"auth_service/internal/configuration"
	"log"
)

func main() {
	config, configError := configuration.ReadConfig("../config.toml")
	if configError != nil {
		log.Fatalf("Error while getting configuration : " + configError.Error())
	}

	app.Run(config)
}
