package entities

type PrintOutResponse struct {
	Route          string            `json:"route"`
	SimpleResponse SimpleResponseMsg `json:"msg"`
}

type SimpleResponseMsg struct {
	Msg string
}

type AuthUseCase interface {
	CreateResponse(inputMsg []byte) string
}

type RedisAuthRepository interface {
}
