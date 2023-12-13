package repository

import (
	"auth_service/internal/domain/models"
	"context"
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

type RedisSessionsRepository struct {
	client *redis.Client
	log    zap.SugaredLogger
}

func (r *RedisSessionsRepository) GetRefreshSession(ctx context.Context, userId string) (*models.RefreshSession, error) {
	return nil, ErrNotImplementedYet
}

func (r *RedisSessionsRepository) InsertRefreshSession(ctx context.Context, session *models.RefreshSession) *error {
	return &ErrNotImplementedYet
}

func NewRedisSessionsRepository(client *redis.Client) *RedisSessionsRepository {
	return &RedisSessionsRepository{client: client}
}
