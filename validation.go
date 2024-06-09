package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type errorBody struct {
	Error string `json:"error"`
}

type reqBody struct {
	Body string `json:"body"`
}

type success struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidateChirpPost(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	res := reqBody{}
	err := decoder.Decode(&res)

	var	errorRes errorBody
	if err != nil {
		errorRes = errorBody{Error: "Something went wrong"}
	}

	if len(res.Body) > 140 {
		errorRes = errorBody{Error: "Chirp is too long"}
	}

	w.Header().Add("Content-Type", "application/json")
	if errorRes.Error != "" {
		data, err := json.Marshal(errorRes)
		if err != nil {
			fmt.Println("couldn't process request")
			return
		}

		w.WriteHeader(http.StatusBadRequest)	
		w.Write(data)
		return
	}

	res.Body = swearWordsCheck(res.Body)

	validRes := success{CleanedBody: res.Body}
	data, err := json.Marshal(validRes)
	if err != nil {
			fmt.Println("couldn't process request")
			return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func swearWordsCheck(body string) string {
	replacement := "****"
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(body)
	for i, _ := range words {
		for _, w := range badWords {
			if strings.ToLower(words[i]) == w {
				words[i] = replacement
			}
		}
	}
	return strings.Join(words, " ")
}