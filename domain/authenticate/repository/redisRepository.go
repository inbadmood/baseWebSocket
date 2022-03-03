package repository

import (
	"BaseWebSocket/entities"
	"context"
	"github.com/go-redis/redis/v8"
)

type redisGamesRepository struct {
	Conn    *redis.Client
	context context.Context
}

// NewRedisRepository 建一個連線
func NewRedisRepository(conn *redis.Client) entities.RedisAuthRepository {
	authContext := context.Background()
	return &redisGamesRepository{
		Conn:    conn,
		context: authContext,
	}
}
