package lib

import (
	"BaseWebSocket/service/server/client"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var apiController *ApiRouter

// 初始化api Router
func ApiControllerInit(newRouter *ApiRouter) {
	apiController = newRouter
}

// 初始化 Client obj
func NewClientObj(hub *client.Hub, clientClientID uint32, conn *websocket.Conn, timeOut time.Duration) (connClient *client.Client) {
	peerSendChannel := make(chan []byte, 2048)
	peerCloseChannel := make(chan []byte, 1)
	peerCloseSendChannel := make(chan []byte, 1)
	connClient = &client.Client{
		Hub:        hub,
		ClientID:   clientClientID,
		Conn:       conn,
		Send:       peerSendChannel,
		Close:      peerCloseChannel,
		CloseSend:  peerCloseSendChannel,
		TimeOutSec: timeOut,
	}

	return connClient
}

// same time to close channel
// 清除離線斷線Client
func disposeClient(p *client.Client) {
	if p != nil {
		p.Hub.UnRegister <- p
		p = nil
	}
}

// readPeerMessage 讀websocket input message
func ReadPeerMessage(p *client.Client) {
	defer func() {
		fmt.Println("readPeerMessage Fail")
		close(p.CloseSend)
		if err := recover(); err != nil {
			fmt.Println("recover from Message Read", fmt.Sprint(err))
		}
	}()

	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			fmt.Println("readPeerMessage Fail")
			return
		}
		// postpone peer timeout
		err = p.Conn.SetReadDeadline(time.Now().Add(p.TimeOutSec * time.Second))
		if err != nil {
			fmt.Println("SetReadDeadline Fail")
		}

		// process input packet
		fmt.Println("readPeerMessage success" + string(message))
		apiController.ProcessPeerMessage(p, message)
	}
}

// writePeerMessage 開啟監聽 Send Channel 寫到client socket
func WritePeerMessage(p *client.Client) {
	defer func() {
		close(p.Close)
		p.Conn.Close()
		fmt.Println("WritePeerMessage Fail")
		disposeClient(p)
	}()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover from Message send")
		}
	}()

	for {
		select {
		case message, ok := <-p.Send:
			if !ok {
				fmt.Println("WritePeerMessage Fail")
				return
			}

			err := p.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("WritePeerMessage Fail" + err.Error())
				return
			}
			// process output packet
			fmt.Println("writePeerMessage success" + string(message))
		case _, ok := <-p.CloseSend:
			if !ok {
				return
			}
		}
	}
}
