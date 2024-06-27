package database

import (
	"errors"
	"os"

	"github.com/sondrefjellving/chirpy/internal/auth"
)

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return User{}, errors.New("user with that email already exists")
		}
	}

	id := len(dbStruct.Chirps) + 1
	user := User{
		Email: email,
		Id: id,
		HashedPassword: hashedPassword,
	}
	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStruct.Users[id]
	if !ok {
		return User{}, os.ErrNotExist 
	}

	user.Email = email
	user.HashedPassword = hashedPassword 
	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UserLogin(email, hashedPassword string) (User, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	hasUser := false
	user := User{}
	for _, currUser := range dbStruct.Users {
		if currUser.Email == email {
			user = currUser
			hasUser = true
			break
		}
	}
	
	if !hasUser {
		return User{}, errors.New("found no user with that email")
	}

	err = auth.CheckPasswordHash(user.HashedPassword, hashedPassword)
	if err != nil {
		return User{}, errors.New("entered password and saved password doesn't match")
	}

	return user, nil
}
