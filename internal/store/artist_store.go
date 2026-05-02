package store

import (
	"database/sql"
	"fmt"
	"sync"

	"OpryScrape/internal/models"
)

type ArtistStore struct {
	mu 				sync.Mutex
	artistMap 		map[string]*models.Artist
	newArtistsCh	chan *models.Artist
}

func NewArtistStore() *ArtistStore {
	return &ArtistStore{
		artistMap: make(map[string]*models.Artist),
		newArtistsCh: make(chan *models.Artist, 100),
	}
}

func (as *ArtistStore) AddArtist(name string) *models.Artist {
	as.mu.Lock()
	defer as.mu.Unlock()

	if artist, exists := as.artistMap[name]; exists{
		return artist
	}

	artist := models.NewArtist(name)
	as.artistMap[name] = artist
	as.newArtistsCh <- artist

	return artist
}


//Used during initial development, before DB was implemented
func (as *ArtistStore) LoopThroughMap() {
	fmt.Println("***Looping through artist map***")
	
	for name, obj := range as.artistMap {
		fmt.Printf("Artist Name: %s Id: %d\n", name, obj.Id)
	}
	fmt.Println("")
}

const artistQuery = `
        SELECT id, name
        FROM artists;
    `
func (as *ArtistStore) LoadFromDB(db *sql.DB) error {
    rows, err := db.Query(artistQuery)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        artist := new(models.Artist)
        if err := rows.Scan(&artist.Id, &artist.Name); err != nil {
            return err
        }

        as.artistMap[artist.Name] = artist
    }

    return nil
}

const artistInsert = `
        INSERT INTO artists (name)
        VALUES (?)
    `
func (as *ArtistStore) InsertArtistsToDb(db *sql.DB, newArtists []*models.Artist) error {
	tx, err := db.Begin()
    if err != nil {
        return err
    }
	defer tx.Rollback()

    stmt, err := tx.Prepare(artistInsert)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, a := range newArtists {
        result, err := stmt.Exec(a.Name)
        if err != nil {
            fmt.Println("Error on Artist name: ", a.Name)
            as.mu.Lock()
            delete(as.artistMap, a.Name)
            as.mu.Unlock()
            continue
        }

		id, err := result.LastInsertId()
		if err != nil {
            fmt.Println("Error getting ID for Artist name: ", a.Name)
            as.mu.Lock()
            delete(as.artistMap, a.Name)
            as.mu.Unlock()
            continue
		}

		a.Id = id
    }

    return tx.Commit()
}