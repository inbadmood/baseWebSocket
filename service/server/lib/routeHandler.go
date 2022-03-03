package lib

import (
	"BaseWebSocket/entities/models"
	"BaseWebSocket/service/server/client"
	"BaseWebSocket/utils"
	"encoding/json"
	"fmt"
	"sync"
)

type ApiRouter struct {
	MiddlewareChain []Middleware
	Handlers        map[string]MessageDispatcher
}

type MessageDispatcher func(client *client.Client, requestRoute string, message []byte) (resp string)

var handlersRWMutex sync.RWMutex

// 注册Dispatcher
func (_r ApiRouter) RegisterMessageDispatcher(key string, value MessageDispatcher) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()

	// 將Middleware包裝上handler
	var mergedHandler = value
	for i := len(_r.MiddlewareChain) - 1; i >= 0; i-- {
		mergedHandler = _r.MiddlewareChain[i](mergedHandler)
	}

	// 註冊
	_r.Handlers[key] = mergedHandler

	return
}

// 取得route Dispatcher
func (_r ApiRouter) GetMessageDispatcher(key string) (value MessageDispatcher, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = _r.Handlers[key]

	return
}

// dispatch inputData
func (_r ApiRouter) ProcessPeerMessage(client *client.Client, message []byte) {
	request := &models.PacketRequestBase{}

	// 檢查input legal
	err := json.Unmarshal(message, request)
	if err != nil {
		errOutput := utils.ErrorMsg(models.ErrJSONMarshal, "Cannot resolve Packet.")
		ResponseToPeerSend(client, []byte(errOutput))
		return
	}

	requestMsg, err := json.Marshal(request.Msg)
	if err != nil {
		errOutput := utils.ErrorMsg(models.ErrJSONMarshal, "Cannot Marshal Packet Msg.")
		ResponseToPeerSend(client, []byte(errOutput))
		return
	}

	inputRoute := request.Route

	// dispatch packet
	// add error handle
	if value, ok := _r.GetMessageDispatcher(inputRoute); ok {
		resp := value(client, inputRoute, requestMsg)
		ResponseToPeerSend(client, []byte(resp))

	} else {
		errOutput := utils.ErrorMsg(models.ErrNoRoute, "")
		ResponseToPeerSend(client, []byte(errOutput))
	}
	return
}

// ResponseToPeerSend 組好的resp 丟進 chan
func ResponseToPeerSend(p *client.Client, msg []byte) {
	if p == nil {
		return
	}

	select {
	case _, ok := <-p.Close:
		fmt.Println("ResponseToPeerSend " + fmt.Sprint(ok))
		return
	default:
		p.Send <- msg
	}
}
