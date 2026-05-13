package store

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/alexnel24/concurrency-opry/internal/db/schema"
	"github.com/alexnel24/concurrency-opry/internal/models"
)

func setupPerformanceTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(schema.PerformancesSchema)
	require.NoError(t, err)
	return db
}

func TestNewPerformanceStore(t *testing.T) {
	ps := NewPerformanceStore()
	assert.NotNil(t, ps.performanceMap)
	assert.NotNil(t, ps.newPerformancesCh)
}

func TestAddPerformance_NewPerformance(t *testing.T) {
	ps := NewPerformanceStore()
	event := models.NewEvent("Opry Live", "https://opry.com/event/1", time.Time{})
	performance := ps.AddPerformance("Brad Paisley", event)

	assert.Equal(t, "Brad Paisley", performance.ArtistName)
	assert.Equal(t, "https://opry.com/event/1", performance.EventLink)
	assert.Len(t, ps.performanceMap, 1)
	assert.Equal(t, int64(1), event.NoOfPerformers)

	received := <-ps.newPerformancesCh
	assert.Equal(t, performance, received)
}

func TestAddPerformance_TwoDistinct(t *testing.T) {
	ps := NewPerformanceStore()
	event := models.NewEvent("Opry Live", "https://opry.com/event/1", time.Time{})
	ps.AddPerformance("Brad Paisley", event)
	ps.AddPerformance("Dolly Parton", event)

	assert.Len(t, ps.performanceMap, 2)
	assert.Len(t, ps.newPerformancesCh, 2)
	assert.Equal(t, int64(2), event.NoOfPerformers)
}

func TestAddPerformance_Duplicate(t *testing.T) {
	ps := NewPerformanceStore()
	event := models.NewEvent("Opry Live", "https://opry.com/event/1", time.Time{})
	first := ps.AddPerformance("Brad Paisley", event)
	second := ps.AddPerformance("Brad Paisley", event)

	assert.Same(t, first, second)
	assert.Len(t, ps.performanceMap, 1)
	assert.Len(t, ps.newPerformancesCh, 1)
	assert.Equal(t, int64(1), event.NoOfPerformers)
}

func TestPerformanceLoadFromDB_EmptyTable(t *testing.T) {
	db := setupPerformanceTestDB(t)
	ps := NewPerformanceStore()

	err := ps.LoadFromDB(db)
	require.NoError(t, err)
	assert.Empty(t, ps.performanceMap)
}

func TestPerformanceLoadFromDB_PopulatesMap(t *testing.T) {
	db := setupPerformanceTestDB(t)

	_, err := db.Exec(`INSERT INTO performances (event_link, artist_name, combo_string) VALUES (?, ?, ?)`,
		"https://opry.com/event/1", "Brad Paisley", "Brad Paisley-https://opry.com/event/1")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO performances (event_link, artist_name, combo_string) VALUES (?, ?, ?)`,
		"https://opry.com/event/1", "Dolly Parton", "Dolly Parton-https://opry.com/event/1")
	require.NoError(t, err)

	ps := NewPerformanceStore()
	err = ps.LoadFromDB(db)
	require.NoError(t, err)

	assert.Len(t, ps.performanceMap, 2)
	assert.Contains(t, ps.performanceMap, "Brad Paisley-https://opry.com/event/1")
	assert.Contains(t, ps.performanceMap, "Dolly Parton-https://opry.com/event/1")
}

func TestInsertPerformancesToDb_AllGood(t *testing.T) {
	db := setupPerformanceTestDB(t)
	ps := NewPerformanceStore()
	event := models.NewEvent("Opry Live", "https://opry.com/event/1", time.Time{})

	performances := []*models.Performance{
		ps.AddPerformance("Brad Paisley", event),
		ps.AddPerformance("Dolly Parton", event),
	}
	<-ps.newPerformancesCh
	<-ps.newPerformancesCh

	err := ps.InsertPerformancesToDb(db, performances)
	require.NoError(t, err)

	assert.NotZero(t, performances[0].Id)
	assert.NotZero(t, performances[1].Id)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM performances").Scan(&count))
	assert.Equal(t, 2, count)
}

func TestInsertPerformancesToDb_OneBadPerformance(t *testing.T) {
	db := setupPerformanceTestDB(t)
	ps := NewPerformanceStore()
	event := models.NewEvent("Opry Live", "https://opry.com/event/1", time.Time{})

	_, err := db.Exec(`INSERT INTO performances (event_link, artist_name, combo_string) VALUES (?, ?, ?)`,
		"https://opry.com/event/1", "Bad Artist", "Bad Artist-https://opry.com/event/1")
	require.NoError(t, err)

	good1 := ps.AddPerformance("Brad Paisley", event)
	bad := ps.AddPerformance("Bad Artist", event)
	good2 := ps.AddPerformance("Dolly Parton", event)
	<-ps.newPerformancesCh
	<-ps.newPerformancesCh
	<-ps.newPerformancesCh

	err = ps.InsertPerformancesToDb(db, []*models.Performance{good1, bad, good2})
	require.NoError(t, err)

	assert.NotZero(t, good1.Id)
	assert.NotZero(t, good2.Id)
	assert.Zero(t, bad.Id)
	assert.NotContains(t, ps.performanceMap, "Bad Artist-https://opry.com/event/1")

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM performances WHERE combo_string != ?", "Bad Artist-https://opry.com/event/1").Scan(&count))
	assert.Equal(t, 2, count)
}

func TestInsertPerformancesToDb_MultipleBadPerformances(t *testing.T) {
	db := setupPerformanceTestDB(t)
	ps := NewPerformanceStore()
	event := models.NewEvent("Opry Live", "https://opry.com/event/1", time.Time{})

	_, err := db.Exec(`INSERT INTO performances (event_link, artist_name, combo_string) VALUES (?, ?, ?)`,
		"https://opry.com/event/1", "Bad Artist One", "Bad Artist One-https://opry.com/event/1")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO performances (event_link, artist_name, combo_string) VALUES (?, ?, ?)`,
		"https://opry.com/event/1", "Bad Artist Two", "Bad Artist Two-https://opry.com/event/1")
	require.NoError(t, err)

	good := ps.AddPerformance("Brad Paisley", event)
	bad1 := ps.AddPerformance("Bad Artist One", event)
	bad2 := ps.AddPerformance("Bad Artist Two", event)
	<-ps.newPerformancesCh
	<-ps.newPerformancesCh
	<-ps.newPerformancesCh

	err = ps.InsertPerformancesToDb(db, []*models.Performance{bad1, good, bad2})
	require.NoError(t, err)

	assert.NotZero(t, good.Id)
	assert.Zero(t, bad1.Id)
	assert.Zero(t, bad2.Id)
	assert.NotContains(t, ps.performanceMap, "Bad Artist One-https://opry.com/event/1")
	assert.NotContains(t, ps.performanceMap, "Bad Artist Two-https://opry.com/event/1")

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM performances WHERE combo_string = ?", "Brad Paisley-https://opry.com/event/1").Scan(&count))
	assert.Equal(t, 1, count)
}

func TestInsertPerformancesToDb_EmptySlice(t *testing.T) {
	db := setupPerformanceTestDB(t)
	ps := NewPerformanceStore()

	err := ps.InsertPerformancesToDb(db, []*models.Performance{})
	require.NoError(t, err)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM performances").Scan(&count))
	assert.Equal(t, 0, count)
}

func TestInsertPerformancesToDb_AllBadPerformances(t *testing.T) {
	db := setupPerformanceTestDB(t)
	ps := NewPerformanceStore()
	event := models.NewEvent("Opry Live", "https://opry.com/event/1", time.Time{})

	_, err := db.Exec(`INSERT INTO performances (event_link, artist_name, combo_string) VALUES (?, ?, ?)`,
		"https://opry.com/event/1", "Bad Artist One", "Bad Artist One-https://opry.com/event/1")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO performances (event_link, artist_name, combo_string) VALUES (?, ?, ?)`,
		"https://opry.com/event/1", "Bad Artist Two", "Bad Artist Two-https://opry.com/event/1")
	require.NoError(t, err)

	bad1 := ps.AddPerformance("Bad Artist One", event)
	bad2 := ps.AddPerformance("Bad Artist Two", event)
	<-ps.newPerformancesCh
	<-ps.newPerformancesCh

	err = ps.InsertPerformancesToDb(db, []*models.Performance{bad1, bad2})
	require.NoError(t, err)

	assert.Zero(t, bad1.Id)
	assert.Zero(t, bad2.Id)
	assert.NotContains(t, ps.performanceMap, "Bad Artist One-https://opry.com/event/1")
	assert.NotContains(t, ps.performanceMap, "Bad Artist Two-https://opry.com/event/1")

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM performances").Scan(&count))
	assert.Equal(t, 2, count)
}
