package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SECONDS_IN_DAY = 3600 * 24
)

func (c *apiConfig) handlerLoginPost(w http.ResponseWriter, req *http.Request) {
	type LoginData struct {
		Password string `json:"password"`
		Email string	`json:"email"`
		Expires_in_seconds int `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(req.Body)
	loginData := LoginData{}
	err := decoder.Decode(&loginData)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request data")
		return 
	}
	
	response, err := c.db.UserLogin(loginData.Email, loginData.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, getJWTClaims(response.Id, loginData.Expires_in_seconds))
	tokenString, err := token.SignedString([]byte(c.jwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating jwt token")
		return
	}
	
	response.Token = tokenString
	respondWithJson(w, http.StatusOK, response)
}


func getJWTClaims(userId, expires_in_seconds int) jwt.Claims {
	if expires_in_seconds > SECONDS_IN_DAY { // set default to a 24h if amount is too high
		expires_in_seconds = SECONDS_IN_DAY
	}

	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(expires_in_seconds))),	
		Subject: strconv.Itoa(userId),
	}

	return claims
}