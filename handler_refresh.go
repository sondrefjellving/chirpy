package main

import (
	"net/http"
	"time"

	"github.com/sondrefjellving/chirpy/internal/auth"
)

const (
	SECONDS_IN_HOUR = 3600
)

func (c *apiConfig) handlerRevokeToken(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = c.db.RevokeRefreshToken(token)
	if err != nil { // no such token in db
		respondWithJson(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}


func (c *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadGateway, err.Error())
		return
	}

	err = c.db.VerifyRefeshToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := c.db.GetUserIdFromRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	hourDuration := time.Second * time.Duration(SECONDS_IN_HOUR)
	accessToken, err := auth.MakeJWT(userId, c.jwtSecret, hourDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "trouble creating refresh token")
		return
	}

	res := response{
		Token: accessToken,
	}

	respondWithJson(w, http.StatusOK, res)
}