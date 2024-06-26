package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sondrefjellving/chirpy/internal/auth"
)

const (
	SECONDS_IN_DAY = 60 * 60 * 24
)

func (c *apiConfig) handlerLoginPost(w http.ResponseWriter, req *http.Request) {
	type LoginData struct {
		Password string `json:"password"`
		Email string	`json:"email"`
		Expires_in_seconds int `json:"expires_in_seconds"`
	}

	type response struct {
		User	
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(req.Body)
	loginData := LoginData{}
	err := decoder.Decode(&loginData)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request data")
		return 
	}

	user, err := c.db.UserLogin(loginData.Email, loginData.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	defaultExpiration := 60 * 60 * 24
	if loginData.Expires_in_seconds == 0 {
		loginData.Expires_in_seconds = defaultExpiration
	} else if loginData.Expires_in_seconds > defaultExpiration {
		loginData.Expires_in_seconds= defaultExpiration
	}

	token, err := auth.MakeJWT(user.Id, c.jwtSecret, time.Duration(loginData.Expires_in_seconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}
	

	respondWithJson(w, http.StatusOK, response{
		User: User{
			Id: user.Id,
			Email: user.Email,
		},
		Token: token,
	}) // here...
}
