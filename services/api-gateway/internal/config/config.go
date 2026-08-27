package config

import (
	"os"
)

type Config struct {
	AuthServiceHost        string
	AuthServicePort        string
	UserServiceHost        string
	UserServicePort        string
	PostServiceHost        string
	PostServicePort        string
	CommentServiceHost     string
	CommentServicePort     string
	LikeServiceHost        string
	LikeServicePort        string
	LeaderboardServiceHost string
	LeaderboardServicePort string
	ServicePort            string
}

func Load() *Config {
	return &Config{
		AuthServiceHost:        getEnv("AUTH_SERVICE_HOST", "auth-service"),
		AuthServicePort:        getEnv("AUTH_SERVICE_PORT", "8081"),
		UserServiceHost:        getEnv("USER_SERVICE_HOST", "user-service"),
		UserServicePort:        getEnv("USER_SERVICE_PORT", "8080"),
		PostServiceHost:        getEnv("POST_SERVICE_HOST", "post-service"),
		PostServicePort:        getEnv("POST_SERVICE_PORT", "8082"),
		CommentServiceHost:     getEnv("COMMENT_SERVICE_HOST", "comment-service"),
		CommentServicePort:     getEnv("COMMENT_SERVICE_PORT", "8083"),
		LikeServiceHost:        getEnv("LIKE_SERVICE_HOST", "like-service"),
		LikeServicePort:        getEnv("Like_SERVICE_PORT", "8084"),
		LeaderboardServiceHost: getEnv("LEADERBOARD_SERVICE_HOST", "leaderboard-service"),
		LeaderboardServicePort: getEnv("LEADERBOARD_SERVICE_PORT", "8085"),
		ServicePort:            getEnv("SERVICE_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
