package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config contains data from .env file.
type Config struct {
	HTTPServerPort string
}

// Load loads .env file to config.
func Load() *Config {
	err := godotenv.Load("../config/config.env")
	if err != nil {
		panic("Failed to load config file: " + err.Error())
	}

	return &Config{
		HTTPServerPort: getEnvRequired("HTTP_SERVER_PORT"),
	}
}

// getEnvRequired extracts data from .env.
func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic("Environment variable " + key + " is required")
	}
	return value
}
