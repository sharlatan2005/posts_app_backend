package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/config"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/httphandler"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/repo/postgres"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/server"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/service"
)

func main() {
	// Прогрузка .env (глобального и локального)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No local .env in cmd/")
	}

	rootEnv := filepath.Join("..", "..", "configs", ".env")
	if err := godotenv.Load(rootEnv); err != nil {
		log.Println("No root .env in configs/")
	}

	cfg := config.Load()
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Connection to database can't be established: %+v", err)
	}

	var userRepo repo.UserRepo
	userRepo = postgres.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := httphandler.NewUserHandler(userService)

	router := server.NewRouter()
	router.SetupRoutes(userHandler)

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
