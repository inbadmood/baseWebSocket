package models

import (
	"BaseWebSocket/process/authentication"
)

type CoreUseCase struct {
	AuthenticUseCase *authentication.UseCase
}
