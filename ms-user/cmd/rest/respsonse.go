package rest

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Info string      `json:"info"`
	Data interface{} `json:"data"`
}

func returnResponse(w http.ResponseWriter, r Response) {
	w.WriteHeader(r.Code)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}
