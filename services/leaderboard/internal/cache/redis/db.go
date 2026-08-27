package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type DB struct {
	Client *redis.Client
}

func NewDB(addr string) (*DB, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := checkConnection(client); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to establish connection to redis: %w", err)
	}

	log.Println("Connection to Redis established!")

	return &DB{Client: client}, nil
}

func checkConnection(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Ping(ctx).Err()
}

func (db *DB) Close() error {
	return db.Client.Close()
}
