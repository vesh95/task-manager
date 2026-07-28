package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func wrireJson(w http.ResponseWriter, resp any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	encErr := json.NewEncoder(w).Encode(resp)
	if encErr != nil {
		log.Printf("error while encoding response: %v", encErr)
	}
}
