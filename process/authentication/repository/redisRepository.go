package repository

import (
	"BaseWebSocket/process/authentication"
	"context"
	"github.com/go-redis/redis/v8"
)

type redisGamesRepository struct {
	Conn    *redis.Client
	context context.Context
}

// NewRedisRepository 建一個連線
func NewRedisRepository(conn *redis.Client) authentication.RedisRepository {
	authContext := context.Background()
	return &redisGamesRepository{
		Conn:    conn,
		context: authContext,
	}
}
