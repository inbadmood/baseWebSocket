package models

type PacketResponseBase struct {
	Code  int         `json:"code"`
	Route string      `json:"route"`
	Msg   interface{} `json:"msg"` // data interface
}

type Pong struct {
	Route string `json:"route"`
}
