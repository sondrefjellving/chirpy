package main

import (
	"encoding/json"
	"net/http"
	"strings"
)


func handlerValidateChirpPost(w http.ResponseWriter, req *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}
	type validRes struct {
		CleanedBody string `json:"cleaned_body"`
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
	res := validRes{CleanedBody: reqData.Body}
	respondWithJson(w, 400, res)
}

func respondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Error decoding message")
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}	

	response := errorResponse{Error: message}
	data, err := json.Marshal(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error decoding error message"))
		return
	}
	w.WriteHeader(status)	
	w.Write(data)
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