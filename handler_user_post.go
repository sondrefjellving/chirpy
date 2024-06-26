package main

import (
	"encoding/json"
	"net/http"

	"github.com/sondrefjellving/chirpy/internal/auth"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (cfg *apiConfig) handlerUserPost(w http.ResponseWriter, req *http.Request) {
	type body struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type response struct {
		User
	}
	decoder := json.NewDecoder(req.Body)
	reqBody := body{}
	err := decoder.Decode(&reqBody)
	w.Header().Add("Text-Content", "application/json")
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, 500, "trouble hashing password")
		return
	}

	user, err := cfg.db.CreateUser(reqBody.Email, hashedPassword)
	if err != nil {
		respondWithError(w, 500, "trouble creating user")
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		User: User{
			Id: user.Id,
			Email: user.Email,
		},
	})
}