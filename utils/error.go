package utils

import (
	"BaseWebSocket/models"
	"encoding/json"
)

func ErrorMsg(code int, Route string) (resp string) {
	res := models.ErrorOutputData{}
	res.Code = code
	res.Route = Route
	json, _ := json.Marshal(res)
	return string(json)
}
