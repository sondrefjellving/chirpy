package main

import (
	"net/http"
	"strconv"

	"github.com/sondrefjellving/chirpy/internal/auth"
)

func (c *apiConfig) handlerChirpDeleteById(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid token")
		return
	}

	userIdString, err := auth.ValidateJWT(token, c.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token doesn't belong to a user")
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't convert string to int")
		return
	}

	chirpId, err := strconv.Atoi(req.PathValue("chirpId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Not valid id")
		return
	}

	err = c.db.CheckIfUserIsOwnerOfChirp(userId, chirpId)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "forbidden action")
		return
	}

	err = c.db.DeleteChirpById(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp")
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}