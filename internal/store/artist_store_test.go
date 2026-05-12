package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/alexnel24/concurrency-opry/internal/db/schema"
	"github.com/alexnel24/concurrency-opry/internal/models"
)

func setupArtistTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(schema.ArtistsSchema)
	require.NoError(t, err)
	return db
}

func TestNewArtistStore(t *testing.T) {
	as := NewArtistStore()
	assert.NotNil(t, as.artistMap)
	assert.NotNil(t, as.newArtistsCh)
}

func TestAddArtist_NewArtist(t *testing.T) {
	as := NewArtistStore()
	artist := as.AddArtist("Brad Paisley")

	assert.Equal(t, "Brad Paisley", artist.Name)
	assert.Len(t, as.artistMap, 1)

	received := <-as.newArtistsCh
	assert.Equal(t, artist, received)
}

func TestAddArtist_TwoDistinctNames(t *testing.T) {
	as := NewArtistStore()
	as.AddArtist("Brad Paisley")
	as.AddArtist("Dolly Parton")

	assert.Len(t, as.artistMap, 2)
	assert.Len(t, as.newArtistsCh, 2)
}

func TestAddArtist_DuplicateName(t *testing.T) {
	as := NewArtistStore()
	first := as.AddArtist("Keith Urban")
	second := as.AddArtist("Keith Urban")

	assert.Same(t, first, second)
	assert.Len(t, as.artistMap, 1)
	assert.Len(t, as.newArtistsCh, 1)
}

func TestArtistLoadFromDB_EmptyTable(t *testing.T) {
	db := setupArtistTestDB(t)
	as := NewArtistStore()

	err := as.LoadFromDB(db)
	require.NoError(t, err)
	assert.Empty(t, as.artistMap)
}

func TestArtistLoadFromDB_PopulatesMap(t *testing.T) {
	db := setupArtistTestDB(t)

	_, err := db.Exec(`INSERT INTO artists (name) VALUES (?)`, "Garth Brooks")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO artists (name) VALUES (?)`, "Luke Combs")
	require.NoError(t, err)

	as := NewArtistStore()
	err = as.LoadFromDB(db)
	require.NoError(t, err)

	assert.Len(t, as.artistMap, 2)
	assert.Equal(t, "Garth Brooks", as.artistMap["Garth Brooks"].Name)
	assert.Equal(t, "Luke Combs", as.artistMap["Luke Combs"].Name)
}

func TestInsertArtistsToDb_AllGood(t *testing.T) {
	db := setupArtistTestDB(t)
	as := NewArtistStore()

	artists := []*models.Artist{
		as.AddArtist("Reba"),
		as.AddArtist("Ella Langley"),
	}
	<-as.newArtistsCh
	<-as.newArtistsCh

	err := as.InsertArtistsToDb(db, artists)
	require.NoError(t, err)

	assert.NotZero(t, artists[0].Id)
	assert.NotZero(t, artists[1].Id)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM artists").Scan(&count))
	assert.Equal(t, 2, count)
}

func TestInsertArtistsToDb_OneBadArtist(t *testing.T) {
	db := setupArtistTestDB(t)
	as := NewArtistStore()

	_, err := db.Exec(`INSERT INTO artists (name) VALUES (?)`, "Bad Artist")
	require.NoError(t, err)

	good1 := as.AddArtist("Paul")
	bad := as.AddArtist("Bad Artist")
	good2 := as.AddArtist("Jimi")
	<-as.newArtistsCh
	<-as.newArtistsCh
	<-as.newArtistsCh

	err = as.InsertArtistsToDb(db, []*models.Artist{good1, bad, good2})
	require.NoError(t, err)

	assert.NotZero(t, good1.Id)
	assert.NotZero(t, good2.Id)
	assert.Zero(t, bad.Id)
	assert.NotContains(t, as.artistMap, "Bad Artist")

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM artists WHERE name != ?", "Bad Artist").Scan(&count))
	assert.Equal(t, 2, count)
}

func TestInsertArtistsToDb_MultipleBadArtists(t *testing.T) {
	db := setupArtistTestDB(t)
	as := NewArtistStore()

	_, err := db.Exec(`INSERT INTO artists (name) VALUES (?)`, "Bad Artist One")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO artists (name) VALUES (?)`, "Bad Artist Two")
	require.NoError(t, err)

	good := as.AddArtist("Carrie Underwood")
	bad1 := as.AddArtist("Bad Artist One")
	bad2 := as.AddArtist("Bad Artist Two")
	<-as.newArtistsCh
	<-as.newArtistsCh
	<-as.newArtistsCh

	err = as.InsertArtistsToDb(db, []*models.Artist{bad1, good, bad2})
	require.NoError(t, err)

	assert.NotZero(t, good.Id)
	assert.Zero(t, bad1.Id)
	assert.Zero(t, bad2.Id)
	assert.NotContains(t, as.artistMap, "Bad Artist One")
	assert.NotContains(t, as.artistMap, "Bad Artist Two")

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM artists WHERE name = ?", "Carrie Underwood").Scan(&count))
	assert.Equal(t, 1, count)
}

func TestInsertArtistsToDb_EmptySlice(t *testing.T) {
	db := setupArtistTestDB(t)
	as := NewArtistStore()

	err := as.InsertArtistsToDb(db, []*models.Artist{})
	require.NoError(t, err)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM artists").Scan(&count))
	assert.Equal(t, 0, count)
}

func TestInsertArtistsToDb_AllBadArtists(t *testing.T) {
	db := setupArtistTestDB(t)
	as := NewArtistStore()

	_, err := db.Exec(`INSERT INTO artists (name) VALUES (?)`, "Bad Artist One")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO artists (name) VALUES (?)`, "Bad Artist Two")
	require.NoError(t, err)

	bad1 := as.AddArtist("Bad Artist One")
	bad2 := as.AddArtist("Bad Artist Two")
	<-as.newArtistsCh
	<-as.newArtistsCh

	err = as.InsertArtistsToDb(db, []*models.Artist{bad1, bad2})
	require.NoError(t, err)

	assert.Zero(t, bad1.Id)
	assert.Zero(t, bad2.Id)
	assert.NotContains(t, as.artistMap, "Bad Artist One")
	assert.NotContains(t, as.artistMap, "Bad Artist Two")

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM artists").Scan(&count))
	assert.Equal(t, 2, count)
}
