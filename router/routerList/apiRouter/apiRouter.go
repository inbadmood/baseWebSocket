package apiRouter

import (
	"BaseWebSocket/models"
	"BaseWebSocket/service/server/client"
	"BaseWebSocket/service/server/lib"
	"BaseWebSocket/utils"
	"encoding/json"
)

type Router struct {
}

func NewRouter(coreUseCase models.CoreUseCase) *lib.ApiRouter {
	newRouterHandler := make(map[string]lib.MessageDispatcher)
	newRouterMiddlewareChain := []lib.Middleware{}
	newRouter := &lib.ApiRouter{
		MiddlewareChain: newRouterMiddlewareChain,
		Handlers:        newRouterHandler,
	}

	// 註註冊middleware LIFO
	lib.UseMiddleware(&newRouter.MiddlewareChain, lib.CheckMessageMiddleware)

	ControllerInit(newRouter)

	return newRouter
}

// 註冊routeDispatcher
func ControllerInit(apiRouter *lib.ApiRouter) {
	apiRouter.RegisterMessageDispatcher("Ping", PingController)
}

// ping route處理
func PingController(client *client.Client, requestRoute string, msg []byte) (resp string) {

	pingRespInterface := &models.Pong{}
	pingRespInterface.Route = "Pong"

	pingResp, err := json.Marshal(pingRespInterface)
	if err != nil {
		errOutput := utils.ErrorMsg(models.ErrJSONMarshal, "Cannot Marshal response Pong.")
		return errOutput
	}
	return string(pingResp)
}
