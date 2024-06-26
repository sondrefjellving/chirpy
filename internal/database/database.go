package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
	RefreshTokens map[string]time.Time
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

func (db *DB) CreateDB() error {
	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp),
		Users: make(map[int]User),
		RefreshTokens: make(map[string]time.Time),
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