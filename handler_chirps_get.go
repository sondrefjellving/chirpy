package main

import (
	"net/http"
	"strconv"

	"github.com/sondrefjellving/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("author_id")	
	sortParam := req.URL.Query().Get("sort")
	w.Header().Add("Content-Type", "application/json")

	chirps := make([]database.Chirp, 0)
	var err error
	if id == "" {
		chirps, err = cfg.db.GetChirps(sortParam)
	} else {
		chirps, err = cfg.db.GetChirpsWithAuthorId(id, sortParam)
	}

	if err != nil {
		respondWithError(w, 500, "Couldn't retrieve chirps")
		return
	}

	respondWithJson(w, 200, chirps)
}

func (cfg *apiConfig) handlerChirpGetById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	dbStruct, err := cfg.db.LoadDB()
	if err != nil {
		respondWithError(w, 500, "Couldn't retrieve chirps")
		return
	}

	id, err := strconv.Atoi(req.PathValue("chirpId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Not valid id")
		return
	}

	chirp, exists := dbStruct.Chirps[id]
	if !exists {
		respondWithError(w, 404, "No chirp with that id")
		return
	}

	respondWithJson(w, 200, chirp)
}