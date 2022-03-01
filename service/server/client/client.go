package client

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub         *Hub
	ClientState int
	ClientID    uint32
	IsClosed    bool
	Conn        *websocket.Conn
	Send        chan []byte
	Close       chan []byte
	CloseSend   chan []byte
	TimeOutSec  time.Duration
}
