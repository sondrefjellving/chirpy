package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)


func (cfg *apiConfig) handlerChirpPost(w http.ResponseWriter, req *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	reqData := reqBody{}
	err := decoder.Decode(&reqData)

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}

	if len(reqData.Body) > 140 {
		respondWithError(w, 500, "Chirp is too long")
	}

	swearWordsCheck(&reqData.Body)
	chirpRes, err := cfg.db.CreateChirp(reqData.Body)
	if err != nil {
		fmt.Println("error creating chirp:", err)
		fmt.Println(chirpRes)
		return
	}

	respondWithJson(w, 201, chirpRes)
}

func swearWordsCheck(body *string) {
	replacement := "****"
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(*body)
	for i := range words {
		for _, w := range badWords {
			if strings.ToLower(words[i]) == w {
				words[i] = replacement
			}
		}
	}
	*body = strings.Join(words, " ")
}