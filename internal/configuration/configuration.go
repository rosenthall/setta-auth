package configuration

import (
	"errors"
	"github.com/spf13/viper"
)

// AuthServiceConfig struct represents configuration of the project
type AuthServiceConfig struct {
	Port                 uint   `mapstructure:"port"`
	PrivateKeyPath       string `mapstructure:"private_key_path"`
	PublicKeyPath        string `mapstructure:"public_key_path"`
	TokenLifeTime        uint   `mapstructure:"token_ttl"`
	RefreshTokenLifeTime uint   `mapstructure:"refresh_token_ttl"`
	LogLevel             string `mapstructure:"log_level"`
	RedisServerIp        string `mapstructure:"redis_address"`
	RedisPassword        string `mapstructure:"redis_password"`
}

// ReadConfig reads configuration file in specified dir and deserializes it
func ReadConfig(path string) (*AuthServiceConfig, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.SetConfigName("AuthService")

	var config AuthServiceConfig
	if err := viper.Sub("AuthService").Unmarshal(&config); err != nil {
		return nil, err
	}

	// Checking if there nil/0/"" fields in config
	if config.Port == 0 ||
		config.PrivateKeyPath == "" ||
		config.PublicKeyPath == "" ||
		config.TokenLifeTime == 0 ||
		config.LogLevel == "" ||
		config.RedisServerIp == "" ||
		config.RedisPassword == "" {
		return nil, errors.New("missing required fields in configuration")
	}

	return &config, nil
}
