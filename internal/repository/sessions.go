package repository

import (
	"auth_service/internal/domain/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	ErrSessionNotExists = errors.New("refresh session doesn't exists")
)

// RefreshSessionsRepository defines the interface for working with refresh sessions and Disconnect to close the connection.
type RefreshSessionsRepository interface {
	GetRefreshSession(ctx context.Context, userId string) (*models.RefreshSession, error)
	InsertRefreshSession(ctx context.Context, session *models.RefreshSession) error
	DeleteRefreshSession(ctx context.Context, userId string) error
	Disconnect() error
}

// RedisRefreshSessionsRepository implements RefreshSessionsRepository using Redis
type RedisRefreshSessionsRepository struct {
	client *redis.Client
	log    zap.SugaredLogger
}

// Disconnect closes redis connection
func (s *RedisRefreshSessionsRepository) Disconnect() error {
	if err := s.client.Close(); err != nil {
		s.log.Errorf("Failed to close Redis client: %v", err)
		return err
	}

	s.log.Info("Successfully closed redis connection!")
	return nil
}

// GetRefreshSession retrieves a refresh session for a given user ID from Redis
func (r *RedisRefreshSessionsRepository) GetRefreshSession(ctx context.Context, userId string) (*models.RefreshSession, error) {
	result, err := r.client.Get(ctx, userId).Result()
	if errors.Is(err, redis.Nil) {
		// Return custom error if session does not exist
		return nil, ErrSessionNotExists
	} else if err != nil {
		// Log and return error if retrieval fails
		r.log.Errorf("Error retrieving refresh session: %v", err)
		return nil, err
	}

	var session models.RefreshSession
	err = json.Unmarshal([]byte(result), &session)
	if err != nil {
		// Log and return error if unmarshaling fails
		r.log.Errorf("Error unmarshaling refresh session: %v", err)
		return nil, err
	}

	return &session, nil
}

// InsertRefreshSession stores a new refresh session in Redis
func (r *RedisRefreshSessionsRepository) InsertRefreshSession(ctx context.Context, session *models.RefreshSession) error {
	sessionData, err := json.Marshal(session)
	if err != nil {
		// Log and return error if marshaling fails
		r.log.Errorf("Error marshaling refresh session: %v", err)
		return err
	}

	_, err = r.client.Set(ctx, session.UserID, sessionData, 0).Result()
	if err != nil {
		// Log and return error if insertion fails
		r.log.Errorf("Error inserting refresh session: %v", err)
		return err
	}

	return nil
}

// DeleteRefreshSession removes a refresh session from Redis
func (r *RedisRefreshSessionsRepository) DeleteRefreshSession(ctx context.Context, userId string) error {
	// Attempt to delete the refresh session from Redis using the user ID
	result, err := r.client.Del(ctx, userId).Result()
	if err != nil {
		// Log and return error if deletion fails
		r.log.Errorf("Error deleting refresh session: %v", err)
		return err
	}

	if result == 0 {
		// If nothing was deleted, it's likely the session didn't exist
		// You can choose to handle this case differently if needed
		return ErrSessionNotExists
	}

	// Return nil if deletion was successful
	return nil
}

// NewRedisRepository creates an instance of RedisRefreshSessionsRepository
func NewRedisRepository(client *redis.Client, logger zap.SugaredLogger) *RedisRefreshSessionsRepository {
	return &RedisRefreshSessionsRepository{
		client: client,
		log:    logger,
	}
}
