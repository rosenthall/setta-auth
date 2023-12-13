package app

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// getRedisConnection creates and returns a new Redis connection.
func getRedisConnection(ctx context.Context, address string, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // ?
	})

	// Checking the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
