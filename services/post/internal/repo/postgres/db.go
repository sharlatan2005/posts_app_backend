package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/config"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(cfg *config.Config) (*DB, error) {
	pool, err := createPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("Pool creation: %w", err)
	}

	err = checkConnection(pool)
	if err != nil {
		return nil, fmt.Errorf("Failed to establish connection to database: %w", err)
	}
	log.Println("Connection to database established!")

	return &DB{Pool: pool}, nil
}

func createPool(cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("Config parse: %w", err)
	}
	poolCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolCfg.MaxConns = 50
	poolCfg.MinConns = 10
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute

	return pgxpool.NewWithConfig(context.Background(), poolCfg)
}

func checkConnection(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pool.Ping(ctx)
}
