package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/events"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/kafka/consumer"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/cache"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/cache/redis"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/config"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/httphandler"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/redishandler"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/repo/postgres"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/server"
	"github.com/sharlatan2005/posts_app_backend/services/leaderboard/internal/service"
)

// Обработчик для лидерборда
func handleActivity(msg []byte) error {
	var act events.Activity
	if err := json.Unmarshal(msg, &act); err != nil {
		return err
	}

	switch act.Type {
	case "like":
		updateScore(act.UserID, 1)
	case "comment":
		updateScore(act.UserID, 5)
	case "post":
		updateScore(act.UserID, 10)
	}
	return nil
}

func updateScore(userID uuid.UUID, points int) {
	fmt.Println(userID, points)
}

func main() {
	brokers := []string{"kafka:9092"}
	consumer, err := consumer.NewMyConsumer(brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	cfg := config.Load()
	pg_db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Connection to PG database can't be established: %+v", err)
	}
	var leaderboardRepo repo.LeaderboardRepo
	leaderboardRepo = postgres.NewLeaderboardRepo(pg_db)

	redis_db, err := redis.NewDB(cfg.RedisAddr)
	if err != nil {
		log.Printf("Connection to Redis database can't be established: %+v", err)
	}
	var leaderboardCache cache.LeaderboardCache
	leaderboardCache = redis.NewLeaderboardCache(redis_db)
	leaderboardService := service.NewLeaderboardService(leaderboardRepo, leaderboardCache)
	activityHandler := redishandler.NewActivityHandler(leaderboardService)

	go func() {
		if err := consumer.Consume("activities", activityHandler.Handle); err != nil {
			log.Printf("Kafka consumer stopped: %v", err)
		}
	}()

	leaderboardHandler := httphandler.NewLeaderboardHandler(leaderboardService)
	router := server.NewRouter()
	router.SetupRoutes(leaderboardHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server running on http://localhost:%s", cfg.ServicePort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
