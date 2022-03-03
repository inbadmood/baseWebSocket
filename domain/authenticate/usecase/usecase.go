package usecase

import (
	"BaseWebSocket/entities"
)

type UseCase struct {
	redisRepo entities.RedisAuthRepository
}

func NewAuthUseCase(r entities.RedisAuthRepository) entities.AuthUseCase {
	return &UseCase{
		redisRepo: r,
	}
}

func (_a UseCase) CreateResponse(inputMsg []byte) string {
	return string(inputMsg)
}
