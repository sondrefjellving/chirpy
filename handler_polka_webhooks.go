package main

import (
	"encoding/json"
	"net/http"
)

const (
	UPGRADED_USER = "user.upgraded"
)

func (c *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, req *http.Request) {
	type params struct {
		Event string `json:"event"`
		Data struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	parameters := params{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse request data")
		return
	}

	if parameters.Event != UPGRADED_USER {
		respondWithJson(w, http.StatusNoContent, struct{}{})
		return
	}

	err = c.db.UpgradeUserToChirpyRed(parameters.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't upgrade user")
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}