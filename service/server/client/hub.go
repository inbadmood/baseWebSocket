package client

import (
	"fmt"
)

type Hub struct {
	IsRunIng       bool
	Clients        map[*Client]bool
	Broadcast      chan []byte
	KickMaintain   chan []byte
	Kick           chan []byte
	kickFromServer chan []byte
	Register       chan *Client
	UnRegister     chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:      make(chan []byte),
		KickMaintain:   make(chan []byte),
		Kick:           make(chan []byte),
		kickFromServer: make(chan []byte),
		Register:       make(chan *Client),
		UnRegister:     make(chan *Client),
		Clients:        make(map[*Client]bool),
	}
}

func (_h *Hub) Run() {
	_h.IsRunIng = true
	defer func() {
		err := recover()
		if err != nil {
			_h.IsRunIng = false
			fmt.Println("Hub Run Error " + fmt.Sprintln(err))
		}
	}()

	for {
		select {
		case client := <-_h.Register:
			_h.Clients[client] = true
		case client := <-_h.UnRegister:
			if _, ok := _h.Clients[client]; ok {
				delete(_h.Clients, client)
			}
		case message := <-_h.Broadcast:
			for client := range _h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(_h.Clients, client)
				}
			}
		case _ = <-_h.kickFromServer:
			for client := range _h.Clients {
				client.Conn.Close()
			}
		case _ = <-_h.KickMaintain:
			for client := range _h.Clients {
				client.Conn.Close()
			}
		}
	}
}

func (_h *Hub) BroadcastToPlayer(msg []byte) {
	_h.Broadcast <- []byte(msg)
}
func (_h *Hub) KickPlayer(msg []byte) {
	_h.Kick <- []byte(msg)
}
func (_h *Hub) ServerMaintain(msg []byte) {
	_h.KickMaintain <- []byte(msg)
}
