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
	ErrSessionNotExists  = errors.New("refresh session doesn't exists")
	ErrNotImplementedYet = errors.New("not implemented yet")
)

type RefreshSessionsRepository interface {
	GetRefreshSession(ctx context.Context, userId string) (*models.RefreshSession, error)
	InsertRefreshSession(ctx context.Context, session *models.RefreshSession) error
}

type RedisRefreshSessionsRepository struct {
	client *redis.Client
	log    zap.SugaredLogger
}

func (r *RedisRefreshSessionsRepository) GetRefreshSession(ctx context.Context, userId string) (*models.RefreshSession, error) {
	result, err := r.client.Get(ctx, userId).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrSessionNotExists
	} else if err != nil {
		r.log.Errorf("Error retrieving refresh session: %v", err)
		return nil, err
	}

	var session models.RefreshSession
	err = json.Unmarshal([]byte(result), &session)
	if err != nil {
		r.log.Errorf("Error unmarshaling refresh session: %v", err)
		return nil, err
	}

	return &session, nil
}

func (r *RedisRefreshSessionsRepository) InsertRefreshSession(ctx context.Context, session *models.RefreshSession) error {
	sessionData, err := json.Marshal(session)
	if err != nil {
		r.log.Errorf("Error marshaling refresh session: %v", err)
		return err
	}

	_, err = r.client.Set(ctx, session.UserID, sessionData, 0).Result()
	if err != nil {
		r.log.Errorf("Error inserting refresh session: %v", err)
		return err
	}

	return nil
}

func NewRedisSessionsRepository(client *redis.Client, logger zap.SugaredLogger) *RedisRefreshSessionsRepository {
	return &RedisRefreshSessionsRepository{
		client: client,
		log:    logger,
	}
}
