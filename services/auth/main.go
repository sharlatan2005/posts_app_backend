package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/user"
	"github.com/sharlatan2005/posts_app_backend/services/auth/internal/config"
	"github.com/sharlatan2005/posts_app_backend/services/auth/internal/httphandler"
	"github.com/sharlatan2005/posts_app_backend/services/auth/internal/server"
	"github.com/sharlatan2005/posts_app_backend/services/auth/internal/service"
)

func main() {
	// Прогрузка .env (глобального и локального)
	// if err := godotenv.Load(".env.local"); err != nil {
	// 	log.Println("No local .env in cmd/")
	// }

	// rootEnv := filepath.Join("..", "..", "configs", ".env")
	// if err := godotenv.Load(rootEnv); err != nil {
	// 	log.Println("No root .env in configs/")
	// }

	cfg := config.Load()

	userClient := user.NewClient(cfg.UserServiceURL)
	userService := service.NewAuthService(cfg.JWTSecret, userClient)
	userHandler := httphandler.NewAuthHandler(userService)

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
