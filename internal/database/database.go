package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`	
}

type DB struct {
	path string
	mux  *sync.Mutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	DB := &DB{
		path: path,
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			return nil, err 
		}
		defer file.Close()
	}

	return DB, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		Body: body,
	}

	id := len(dbStruct.Chirps)
	for {
		_, exists := dbStruct.Chirps[id]
		if !exists {
			chirp.Id = id
			dbStruct.Chirps[id] = chirp 
		}
	}

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err 
	}
	return chirp, nil
}

func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(db.path)
		if err != nil {
			return err 
		}
		defer file.Close()
	}
	return nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, len(dbData.Chirps)) 
	i := 0
	for _, chirp := range dbData.Chirps {
		chirps[i] = chirp
		i++
	}

	return chirps, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	if err := db.ensureDB(); err != nil {
		return DBStructure{}, err 
	}

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbStruct := DBStructure{}
	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStruct, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	file, err := os.Create(db.path)
	if err != nil {
		return err
	}
	defer file.Close()

	dbAsJson, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	_, err = file.Write(dbAsJson)
	if err != nil {
		return err 
	}
	return nil
}