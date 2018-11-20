package helper

import (
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   error       `json:"error"`
}

func Response(w http.ResponseWriter, code int, Status bool, Message string, Data interface{}, Error error) {
	response := ResponseData{
		Status,
		Message,
		Data,
		Error,
	}

	rj, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(rj)
}
