package database

import (
	"errors"
	"time"

	"github.com/sondrefjellving/chirpy/internal/auth"
)

const (
	REFRESH_TOKEN_DURATION = 60*60*24*60 // 60 days in seconds, 60 seconds * 60 minutes * 24h in a day * 60 days
)

func (db *DB) GetRefreshToken(userId int) (string, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateRefreshToken()
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

func (db *DB) VerifyRefeshToken(token string) error {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return err
	}

	if _, ok := dbStruct.RefreshTokens[token]; !ok {
		return errors.New("invalid token")
	}
	return nil
}

func (db *DB) AddRefreshToken(token string, durationInSeconds int) error {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(time.Duration(time.Duration(durationInSeconds).Seconds()))
	dbStruct.RefreshTokens[token] = expiresAt
	if err := db.writeDB(dbStruct); err != nil {
		return err
	}

	return nil
}