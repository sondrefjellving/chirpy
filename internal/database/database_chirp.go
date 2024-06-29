package database

import (
	"errors"
	"sort"
	"strconv"
)

type Chirp struct {
	Id int `json:"id"`
	AuthorId int `json:"author_id"`
	Body string `json:"body"`	
}

func (db *DB) CheckIfUserIsOwnerOfChirp(userId, chirpId int) error {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return err
	}
	
	chirp, ok := dbStruct.Chirps[chirpId]
	if !ok {
		return errors.New("couldn't find chirp with that id")
	}

	if chirp.AuthorId != userId {
		return errors.New("user is not authorized to delete chirp with that id")
	}
	return nil
}

func (db *DB) DeleteChirpById(id int) error {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return err
	}

	if _, ok := dbStruct.Chirps[id]; !ok {
		return errors.New("no chirp with that id")
	}

	delete(dbStruct.Chirps, id)
	if err := db.writeDB(dbStruct); err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateChirp(userId int, body string) (Chirp, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := 1
	chirp := Chirp{
		Id: id,
		AuthorId: userId,
		Body: body,
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

func (db *DB) GetChirpsWithAuthorId(authorId string) ([]Chirp, error) {
	id, err := strconv.Atoi(authorId)
	if err != nil {
		return nil, err
	}

	dbData, err := db.LoadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbData.Chirps)) 
	for _, chirp := range dbData.Chirps {
		if chirp.AuthorId == id {
			chirps = append(chirps, chirp)
		}
	}

	return sortChirps(chirps), nil
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

	return sortChirps(chirps), nil
}

func sortChirps(chirps []Chirp) []Chirp {
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	return chirps
}