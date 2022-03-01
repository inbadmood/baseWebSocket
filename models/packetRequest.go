package models

// add request code string map

type PacketRequestBase struct {
	Route string      `json:"route"` // packet Type
	Msg   interface{} `json:"msg"`   // data interface
}
