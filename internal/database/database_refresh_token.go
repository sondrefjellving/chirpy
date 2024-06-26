package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

const (
	REFRESH_TOKEN_DURATION = 60*60*24*60 // 60 days in seconds, 60 seconds * 60 minutes * 24h in a day * 60 days
)

func (db *DB) GetRefreshToken(userId int) (string, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return "", err
	}

	token, err := createRefreshToken()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(time.Duration(REFRESH_TOKEN_DURATION * time.Second))
	dbStruct.RefreshTokens[token] = expiresAt

	user, ok := dbStruct.Users[userId]
	if !ok {
		return "", errors.New("trouble finding user from db")
	}

	user.RefreshToken = token
	dbStruct.Users[userId] = user

	if err := db.writeDB(dbStruct); err != nil {
		return "", err
	}
	return token, nil
}

func createRefreshToken() (refreshToken string, err error) {
	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", errors.New("couldn't get random string for refresh token") 
	}

	refreshToken = hex.EncodeToString(randomBytes)
	return refreshToken, nil
}