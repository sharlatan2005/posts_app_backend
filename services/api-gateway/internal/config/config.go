package config

import (
	"os"
)

type Config struct {
	AuthServiceHost string
	AuthServicePort string
	UserServiceHost string
	UserServicePort string
	ServicePort     string
}

func Load() *Config {
	return &Config{
		AuthServiceHost: getEnv("AUTH_SERVICE_HOST", "auth-service"),
		AuthServicePort: getEnv("AUTH_SERVICE_PORT", "8080"),
		UserServiceHost: getEnv("USER_SERVICE_HOST", "user-service"),
		UserServicePort: getEnv("USER_SERVICE_PORT", "8080"),
		ServicePort:     getEnv("SERVICE_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
