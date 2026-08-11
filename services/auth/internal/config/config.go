package config

import "os"

type Config struct {
	JWTSecret      string
	UserServiceURL string
	ServicePort    string
}

func Load() *Config {
	return &Config{
		JWTSecret:      getEnv("JWT_SECRET", "localhost"),
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8080"),
		ServicePort:    getEnv("SERVICE_PORT", "8081"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
