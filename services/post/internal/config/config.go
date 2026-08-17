package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	ServicePort    string
	JWTSecret      string
	UserServiceURL string
}

func Load() *Config {
	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "posts_app_db"),
		DBSSLMode:      getEnv("DB_SSL_MODE", "disable"),
		ServicePort:    getEnv("SERVICE_PORT", "8082"),
		JWTSecret:      getEnv("JWT_SECRET", "my_beautiful_secret"),
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://user-service:8080"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
