package database

type Chirp struct {
	Id int `json:"id"`
	AuthorId int `json:"author_id"`
	Body string `json:"body"`	
}

func (db *DB) CreateChirp(userId int, body string) (Chirp, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStruct.Chirps) + 1
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