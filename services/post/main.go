package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/auth"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/config"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/httphandler"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/repo/postgres"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/server"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/service"
)

func main() {
	// Прогрузка .env (глобального и локального)

	// if err := godotenv.Load(".env.local"); err != nil {
	// 	log.Println("No local .env in cmd/")
	// }

	// rootEnv := filepath.Join("..", "..", "configs", ".env.global")
	// if err := godotenv.Load(rootEnv); err != nil {
	// 	log.Println("No root .env in configs/")
	// }

	cfg := config.Load()
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Connection to database can't be established: %+v", err)
	}

	var postRepo repo.PostRepo
	postRepo = postgres.NewPostRepo(db)
	postService := service.NewPostService(postRepo)
	postHandler := httphandler.NewPostHandler(postService)

	router := server.NewRouter()
	authStruct := auth.NewAuth([]byte(cfg.JWTSecret))
	router.SetupRoutes(authStruct, postHandler)

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
