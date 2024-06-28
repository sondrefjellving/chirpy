package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/sondrefjellving/chirpy/internal/auth"
)


func (c *apiConfig) handlerChirpPost(w http.ResponseWriter, req *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user not logged in")
		return
	}

	userIdString, err := auth.ValidateJWT(token, c.jwtSecret)	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error validating token")
		return
	}

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error converting string to int")
		return
	}

	decoder := json.NewDecoder(req.Body)
	reqData := reqBody{}
	err = decoder.Decode(&reqData)

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(reqData.Body) > 140 {
		respondWithError(w, 500, "Chirp is too long")
		return
	}

	cleanedBody := getCleanBody(reqData.Body)
	chirpRes, err := c.db.CreateChirp(userId, cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, 201, chirpRes)
}

func getCleanBody(body string) string {
	replacement := "****"
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(body)
	for i := range words {
		for _, w := range badWords {
			if strings.ToLower(words[i]) == w {
				words[i] = replacement
			}
		}
	}
	return strings.Join(words, " ")
}