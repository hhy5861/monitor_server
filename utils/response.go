package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Response struct {
		Code    int         `json:"code"`
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
)

func (res *Response) WriteJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	body, _ := json.Marshal(res)
	fmt.Fprint(w, string(body))
}
