package main

import (
	"net/http"

	"github.com/sondrefjellving/chirpy/internal/auth"
)

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
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "trouble creating refresh token")
		return
	}

	err = c.db.AddRefreshToken(refreshToken, 3600)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "trouble adding refresh token to db")
		return
	}

	res := response{
		Token: refreshToken,
	}

	respondWithJson(w, http.StatusOK, res)
}