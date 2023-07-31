package response

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, msg string, status int, data, meta any) {
	var j struct {
		Msg  string `json:"msg"`
		Data any    `json:"data"`
		Meta any    `json:"meta"`
	}

	j.Data = data
	j.Meta = meta
	j.Msg = msg

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(j)
}

func Message(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"msg"`
	}

	j.Msg = msg
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(j)
}

func Error(w http.ResponseWriter, status int, err error) {
	Message(w, status, err.Error())
}
