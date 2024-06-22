package main

import (
	"encoding/json"
	"net/http"

	"github.com/sondrefjellving/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserPost(w http.ResponseWriter, req *http.Request) {
	type body struct {
		Password string `json:"password"`
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

	user, err := cfg.db.CreateUser(reqBody.Email, reqBody.Password)
	if err != nil {
		respondWithError(w, 500, "Trouble creating user")
		return
	}

	respondWithJson(w, 201, database.UserDTO{
		Id: user.Id,
		Email: user.Email,
	})
}