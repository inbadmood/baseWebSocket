package usecase

import "BaseWebSocket/process/authentication"

type UseCase struct {
	redisRepo authentication.RedisRepository
}

func NewAuthUseCase(r authentication.RedisRepository) authentication.UseCase {
	return &UseCase{
		redisRepo: r,
	}
}
