package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/config"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/domain"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/repo"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/repo/postgres"
)

func main() {
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
	user := &domain.User{
		ID:            uuid.New(),
		Username:      "pidor",
		Password_hash: "sdas",
		Name:          "Pidor",
		Surname:       "Pidorovich",
	}

	err = userRepo.Create(context.Background(), user)

	if err != nil {
		log.Fatalf("User creation: %+v", err)
	}

	fmt.Printf("Score: %d, created_at: %s", user.Score, user.Created_at.String())
}
