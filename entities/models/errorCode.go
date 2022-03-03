package models

const (
	ErrJSONMarshal = 10001 // 轉格式json失敗
	ErrNoRoute     = 10002 // 無相關route
	ErrTimeOut     = 10003 // 閒置斷線
)

type ErrorOutputData struct {
	Code  int    `json:"code"`
	Route string `json:"route"` // error Type
}
