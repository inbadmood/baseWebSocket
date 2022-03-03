package delivery

import (
	"BaseWebSocket/entities"
	"BaseWebSocket/entities/models"
	"BaseWebSocket/service/server/client"
	"BaseWebSocket/service/server/lib"
	"BaseWebSocket/utils"
	"encoding/json"
)

type handler struct {
	authUseCase entities.AuthUseCase
}

var funcHandler *handler

func NewRouter(newAuthUseCase entities.AuthUseCase) *lib.ApiRouter {
	newRouterHandler := make(map[string]lib.MessageDispatcher)
	newRouterMiddlewareChain := []lib.Middleware{}
	newRouter := &lib.ApiRouter{
		MiddlewareChain: newRouterMiddlewareChain,
		Handlers:        newRouterHandler,
	}
	funcHandler = &handler{
		authUseCase: newAuthUseCase,
	}

	// 註註冊middleware LIFO
	lib.UseMiddleware(&newRouter.MiddlewareChain, lib.CheckMessageMiddleware)

	ControllerInit(newRouter)

	return newRouter
}

// 註冊routeDispatcher
func ControllerInit(apiRouter *lib.ApiRouter) {
	apiRouter.RegisterMessageDispatcher("Ping", PingController)
	apiRouter.RegisterMessageDispatcher("PrintOut", PrintOutController)
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
func PrintOutController(client *client.Client, requestRoute string, msg []byte) (resp string) {
	printOutRespInterface := &entities.PrintOutResponse{
		Route: requestRoute,
	}
	printOutMsg := funcHandler.authUseCase.CreateResponse(msg)
	printOutRespInterface.SimpleResponse.Msg = printOutMsg

	printOutResp, err := json.Marshal(printOutRespInterface)
	if err != nil {
		errOutput := utils.ErrorMsg(models.ErrJSONMarshal, "Cannot Marshal response PrintOut.")
		return errOutput
	}
	return string(printOutResp)
}
