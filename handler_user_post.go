package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUserPost(w http.ResponseWriter, req *http.Request) {
	type body struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	reqBody := body{}
	err := decoder.Decode(&reqBody)
	w.Header().Add("Text-Content", "application/json")
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(reqBody.Email)
	if err != nil {
		respondWithError(w, 500, "Trouble creating user")
		return
	}

	respondWithJson(w, 201, user)
}