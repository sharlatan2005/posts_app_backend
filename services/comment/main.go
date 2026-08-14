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
	"github.com/sharlatan2005/chat_app_go_backend_pkg/auth"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/config"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/httphandler"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/repo/postgres"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/server"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/service"
)

func main() {
	// Прогрузка .env (глобального и локального)

	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No local .env in cmd/")
	}

	rootEnv := filepath.Join("..", "..", "configs", ".env.global")
	if err := godotenv.Load(rootEnv); err != nil {
		log.Println("No root .env in configs/")
	}

	cfg := config.Load()
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Connection to database can't be established: %+v", err)
	}

	var commentRepo repo.CommentRepo
	commentRepo = postgres.NewCommentRepo(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := httphandler.NewCommentHandler(commentService)

	router := server.NewRouter()
	authStruct := auth.NewAuth([]byte(cfg.JWTSecret))
	router.SetupRoutes(authStruct, commentHandler)

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
