package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sondrefjellving/chirpy/internal/auth"
	"github.com/sondrefjellving/chirpy/internal/database"
)

func (c *apiConfig) handlerUserPut(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		database.UserDTO
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, c.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate jwt")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse user id")
	}

	user, err := c.db.UpdateUser(userIDInt, params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update user")
	}

	respondWithJson(w, http.StatusOK, response{
		UserDTO: database.UserDTO{
			Id: user.Id,
			Email: user.Email,
		},
	})
}