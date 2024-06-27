package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sondrefjellving/chirpy/internal/auth"
)

const (
	SECONDS_IN_60_DAYS = 60 * 60 * 24 * 60
)

func (c *apiConfig) handlerLoginPost(w http.ResponseWriter, req *http.Request) {
	type LoginData struct {
		Password string	`json:"password"`
		Email string	`json:"email"`
	}

	type response struct {
		User	
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	accessToken, err := 
		auth.MakeJWT(user.Id, c.jwtSecret, time.Duration(SECONDS_IN_HOUR*time.Second))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token")
		return
	}

	err = c.db.SaveRefreshToken(user.Id, refreshToken, SECONDS_IN_60_DAYS)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "trouble saving refresh token to db")
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			Id: user.Id,
			Email: user.Email,
		},
		Token: accessToken,
		RefreshToken: refreshToken,
	}) // here...
}
