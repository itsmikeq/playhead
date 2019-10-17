package app

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// A secret string used for session cookies, passwords, etc.
	SecretKey []byte
	JwtKey    []byte
}

func InitConfig() (*Config, error) {
	config := &Config{
		SecretKey: []byte(viper.GetString("secret_key")),
		JwtKey:    []byte(viper.GetString("jwt_key")),
	}
	if len(config.SecretKey) == 0 {
		return nil, fmt.Errorf("SecretKey must be set")
	}
	if len(config.JwtKey) == 0 {
		return nil, fmt.Errorf("JwtKey must be set")
	}
	return config, nil
}
