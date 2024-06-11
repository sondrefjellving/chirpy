package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps()
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		respondWithError(w, 500, "Couldn't retrieve chirps")
	}

	chirpsAsJson, err := json.Marshal(chirps)
	if err != nil {
		respondWithError(w, 500, "Couldn't convert chirps to json format")
	}

	w.Write(chirpsAsJson)
}