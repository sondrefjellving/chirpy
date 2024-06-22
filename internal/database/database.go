package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`	
}

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`	
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

func NewDB(path string, debugMode *bool) (*DB, error) {
	db := &DB{
		path: path,
		mux: &sync.RWMutex{},
	}

	if *debugMode {
		return db, db.CreateDB()
	}

	err := db.ensureDB() 
	return db, err 
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStruct.Chirps) + 1
	user := User{
		Email: email,
		Id: id,
	}

	for {
		_, exists := dbStruct.Users[id]
		if !exists {
			user.Id = id
			dbStruct.Users[id] = user 
			break
		}
		id++
	}

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err 
	}
	return user, nil

}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{
		Body: body,
		Id: id,
	}

	for {
		_, exists := dbStruct.Chirps[id]
		if !exists {
			chirp.Id = id
			dbStruct.Chirps[id] = chirp 
			break
		}
		id++
	}

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err 
	}
	return chirp, nil
}

func (db *DB) CreateDB() error {
	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp),
		Users: make(map[int]User),
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.CreateDB()
	}
	return nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbData, err := db.LoadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbData.Chirps)) 
	for _, chirp := range dbData.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) LoadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	
	dbStructure := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbAsJson, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dbAsJson, 0600)
	if err != nil {
		return err 
	}
	return nil
}