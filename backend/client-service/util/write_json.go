package util

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Code int `json:"status"`
	Msg  any `json:"message"`
}

func WriteJson(w http.ResponseWriter, code int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(&response{Code: code, Msg: msg})
}
