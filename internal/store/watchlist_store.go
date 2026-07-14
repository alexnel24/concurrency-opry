package store

import (
	"database/sql"

	"github.com/alexnel24/concurrency-opry/internal/models"
)

type WatchlistStore struct{}

func NewWatchlistStore() *WatchlistStore {
	return &WatchlistStore{}
}

const watchedArtistsQuery = `
        SELECT id, name, artist_id, owner_id
        FROM watched_artists
        WHERE owner_id IS ?;
    `

func (ws *WatchlistStore) GetWatchedArtists(db *sql.DB, ownerId *int64) ([]*models.WatchedArtist, error) {
	rows, err := db.Query(watchedArtistsQuery, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	watched := make([]*models.WatchedArtist, 0)
	for rows.Next() {
		wa := new(models.WatchedArtist)
		var artistId, ownerIdCol sql.NullInt64
		if err := rows.Scan(&wa.Id, &wa.Name, &artistId, &ownerIdCol); err != nil {
			return nil, err
		}

		if artistId.Valid {
			wa.ArtistId = &artistId.Int64
		}
		if ownerIdCol.Valid {
			wa.OwnerId = &ownerIdCol.Int64
		}

		watched = append(watched, wa)
	}

	return watched, nil
}

const watchedArtistLookup = `
        SELECT id FROM watched_artists WHERE name = ? AND owner_id IS ?;
    `

const watchedArtistInsert = `
        INSERT INTO watched_artists (name, artist_id, owner_id)
        VALUES (?, ?, ?);
    `

func (ws *WatchlistStore) AddWatchedArtist(db *sql.DB, wa *models.WatchedArtist) (alreadyExists bool, err error) {
	var existingId int64
	err = db.QueryRow(watchedArtistLookup, wa.Name, wa.OwnerId).Scan(&existingId)
	if err == nil {
		wa.Id = existingId
		return true, nil
	}
	if err != sql.ErrNoRows {
		return false, err
	}

	result, err := db.Exec(watchedArtistInsert, wa.Name, wa.ArtistId, wa.OwnerId)
	if err != nil {
		return false, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return false, err
	}

	wa.Id = id
	return false, nil
}

const watchedArtistDelete = `
        DELETE FROM watched_artists WHERE name = ? AND owner_id IS ?;
    `

func (ws *WatchlistStore) RemoveWatchedArtist(db *sql.DB, name string, ownerId *int64) (found bool, err error) {
	result, err := db.Exec(watchedArtistDelete, name, ownerId)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
