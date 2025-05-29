package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, err string) {
	if code > 499 {
		log.Printf("Responding with error: %v", err)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJson(w, code, errorResponse{
		Error: err,
	})
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling Json: %s", err)
	}
	w.WriteHeader(code)
	w.Write(data)
}
