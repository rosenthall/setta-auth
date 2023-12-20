package app

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

// getRedisConnection creates and returns a new Redis connection.
func getRedisConnection(ctx context.Context, address string, password string, log *zap.SugaredLogger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 1,
		IdleTimeout:  1 * time.Hour,
	})

	// Checking the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Errorf("Error while ping Redis server: %v", err)
		return nil, err
	}

	return client, nil
}
