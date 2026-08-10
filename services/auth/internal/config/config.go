package config

import "os"

type Config struct {
	JWTSecret string
}

func Load() *Config {
	return &Config{
		JWTSecret: getEnv("JWT_SECRET", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
