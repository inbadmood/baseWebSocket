package router

import (
	"BaseWebSocket/models"
	"BaseWebSocket/router/routerList/apiRouter"
	"BaseWebSocket/service/server/lib"
)

func NewApiRouter(coreUseCase models.CoreUseCase) *lib.ApiRouter {
	var newRouters *lib.ApiRouter
	newRouters = apiRouter.NewRouter(coreUseCase)
	return newRouters
}
