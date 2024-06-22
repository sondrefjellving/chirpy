package main

import (
	"encoding/json"
	"net/http"
)

func (c *apiConfig) handlerLoginPost(w http.ResponseWriter, req *http.Request) {
	type LoginData struct {
		Password string `json:"password"`
		Email string	`json:"email"`
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
	
	respondWithJson(w, http.StatusOK, response)
}