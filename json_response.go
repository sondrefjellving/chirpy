package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Error decoding message")
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}	

	response := errorResponse{Error: message}
	data, err := json.Marshal(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error decoding error message"))
		return
	}
	w.WriteHeader(status)	
	w.Write(data)
}
