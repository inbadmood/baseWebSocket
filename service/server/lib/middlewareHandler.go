package lib

import (
	"BaseWebSocket/entities/models"
	"BaseWebSocket/service/server/client"
	"BaseWebSocket/utils"
	"encoding/json"
)

// Middleware is public middleware
type Middleware func(MessageDispatcher) MessageDispatcher

// Use add middleware
func UseMiddleware(middlewareChain *[]Middleware, middleware Middleware) {
	*middlewareChain = append(*middlewareChain, middleware)
}

func CheckMessageMiddleware(next MessageDispatcher) MessageDispatcher {
	return func(client *client.Client, requestRoute string, message []byte) (resp string) {
		request := &models.PacketRequestBase{}

		err := json.Unmarshal(message, request)
		if err != nil {
			errOutput := utils.ErrorMsg(models.ErrJSONMarshal, "Cannot resolve Packet.")
			ResponseToPeerSend(client, []byte(errOutput))
			return
		}

		_, err = json.Marshal(request.Msg)
		if err != nil {
			errOutput := utils.ErrorMsg(models.ErrJSONMarshal, "Cannot Marshal Packet Msg.")
			ResponseToPeerSend(client, []byte(errOutput))
			return
		}

		resp = next(client, requestRoute, message)
		return resp
	}
}
