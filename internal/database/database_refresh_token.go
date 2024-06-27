package database

import (
	"errors"
	"time"
)

const (
	REFRESH_TOKEN_DURATION = 60*60*24*60 // 60 days in seconds, 60 seconds * 60 minutes * 24h in a day * 60 days
)

type RefreshToken struct {
	UserId		int
	ExpiresAt 	time.Time
	Token		string
}

func (db *DB) GetUserIdFromRefreshToken(token string) (int, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return 0, err
	}

	refreshToken, ok := dbStruct.RefreshTokens[token]
	if !ok {
		return 0, errors.New("invalid refresh token")
	}

	return refreshToken.UserId, nil
}

func (db *DB) SaveRefreshToken(userId int, token string, expiresInSeconds int) error {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(expiresInSeconds))
	dbStruct.RefreshTokens[token] = RefreshToken{
		UserId: userId,
		ExpiresAt: expiresAt,
		Token: token,
	}

	if err := db.writeDB(dbStruct); err != nil {
		return err
	}

	return nil
}

func (db *DB) VerifyRefeshToken(token string) error {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return err
	}
	entry, ok := dbStruct.RefreshTokens[token]
	if !ok {
		return errors.New("invalid token")
	}

	if entry.ExpiresAt.Before(time.Now()) {
		return errors.New("token has expired")
	}
	
	return nil
}

func (db *DB) RevokeRefreshToken(token string) error {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return err
	}

	if _, ok := dbStruct.RefreshTokens[token]; !ok {
		return errors.New("invalid token")	
	}

	delete(dbStruct.RefreshTokens, token)

	err = db.writeDB(dbStruct)
	if err != nil {
		return errors.New("touble writing to db")
	}

	return nil
}