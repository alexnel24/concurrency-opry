package store

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"OpryScrape/internal/db/schema"
	"OpryScrape/internal/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(schema.EventsSchema)
	require.NoError(t, err)
	return db
}

func TestNewEventStore(t *testing.T) {
	es := NewEventStore()
	assert.NotNil(t, es.EventMap)
	assert.NotNil(t, es.newEventsCh)
}

func TestAddEvent_NewEvent(t *testing.T) {
	es := NewEventStore()
	event := es.AddEvent("Show Title", "https://opry.com/event/1")

	assert.Equal(t, "Show Title", event.Title)
	assert.Equal(t, "https://opry.com/event/1", event.Link)
	assert.Len(t, es.EventMap, 1)

	received := <-es.newEventsCh
	assert.Equal(t, event, received)
}

func TestAddEvent_TwoDistinctLinks(t *testing.T) {
	es := NewEventStore()
	es.AddEvent("Show One", "https://opry.com/event/1")
	es.AddEvent("Show Two", "https://opry.com/event/2")

	assert.Len(t, es.EventMap, 2)
	assert.Len(t, es.newEventsCh, 2)
}

func TestAddEvent_DuplicateLink(t *testing.T) {
	es := NewEventStore()
	first := es.AddEvent("Show Title", "https://opry.com/event/1")
	second := es.AddEvent("Show Title", "https://opry.com/event/1")

	assert.Same(t, first, second)
	assert.Len(t, es.EventMap, 1)
	assert.Len(t, es.newEventsCh, 1)
}

func TestLoadFromDB_EmptyTable(t *testing.T) {
	db := setupTestDB(t)
	es := NewEventStore()

	err := es.LoadFromDB(db)
	require.NoError(t, err)
	assert.Empty(t, es.EventMap)
}

func TestLoadFromDB_PopulatesMap(t *testing.T) {
	db := setupTestDB(t)
	now := time.Now().UTC().Truncate(time.Second)

	_, err := db.Exec(`INSERT INTO events (link, title, time, no_of_performers) VALUES (?, ?, ?, ?)`,
		"https://opry.com/event/1", "Show One", now.Format(time.RFC3339), 0)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO events (link, title, time, no_of_performers) VALUES (?, ?, ?, ?)`,
		"https://opry.com/event/2", "Show Two", now.Format(time.RFC3339), 2)
	require.NoError(t, err)

	es := NewEventStore()
	err = es.LoadFromDB(db)
	require.NoError(t, err)

	assert.Len(t, es.EventMap, 2)
	assert.Equal(t, "Show One", es.EventMap["https://opry.com/event/1"].Title)
	assert.Equal(t, "Show Two", es.EventMap["https://opry.com/event/2"].Title)
}

func TestInsertEventsToDb_AllGood(t *testing.T) {
	db := setupTestDB(t)
	es := NewEventStore()

	events := []*models.Event{
		es.AddEvent("Show One", "https://opry.com/event/1"),
		es.AddEvent("Show Two", "https://opry.com/event/2"),
	}
	// drain channel
	<-es.newEventsCh
	<-es.newEventsCh

	err := es.InsertEventsToDb(db, events)
	require.NoError(t, err)

	assert.NotZero(t, events[0].Id)
	assert.NotZero(t, events[1].Id)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count))
	assert.Equal(t, 2, count)
}

func TestInsertEventsToDb_OneBadEvent(t *testing.T) {
	db := setupTestDB(t)
	es := NewEventStore()

	// Pre-insert to make "bad" a duplicate link that will violate UNIQUE constraint
	_, err := db.Exec(`INSERT INTO events (link, title, time, no_of_performers) VALUES (?, ?, datetime('now'), 0)`,
		"https://opry.com/event/bad", "Pre-existing")
	require.NoError(t, err)

	good1 := es.AddEvent("Good One", "https://opry.com/event/good1")
	bad := es.AddEvent("Bad Event", "https://opry.com/event/bad")
	good2 := es.AddEvent("Good Two", "https://opry.com/event/good2")
	<-es.newEventsCh
	<-es.newEventsCh
	<-es.newEventsCh

	err = es.InsertEventsToDb(db, []*models.Event{good1, bad, good2})
	require.NoError(t, err)

	assert.NotZero(t, good1.Id)
	assert.NotZero(t, good2.Id)
	assert.Zero(t, bad.Id)
	assert.NotContains(t, es.EventMap, "https://opry.com/event/bad")

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM events WHERE link != ?", "https://opry.com/event/bad").Scan(&count))
	assert.Equal(t, 2, count)
}

func TestInsertEventsToDb_MultipleBadEvents(t *testing.T) {
	db := setupTestDB(t)
	es := NewEventStore()

	_, err := db.Exec(`INSERT INTO events (link, title, time, no_of_performers) VALUES (?, ?, datetime('now'), 0)`,
		"https://opry.com/event/bad1", "Pre-existing One")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO events (link, title, time, no_of_performers) VALUES (?, ?, datetime('now'), 0)`,
		"https://opry.com/event/bad2", "Pre-existing Two")
	require.NoError(t, err)

	good := es.AddEvent("Good Event", "https://opry.com/event/good")
	bad1 := es.AddEvent("Bad One", "https://opry.com/event/bad1")
	bad2 := es.AddEvent("Bad Two", "https://opry.com/event/bad2")
	<-es.newEventsCh
	<-es.newEventsCh
	<-es.newEventsCh

	err = es.InsertEventsToDb(db, []*models.Event{bad1, good, bad2})
	require.NoError(t, err)

	assert.NotZero(t, good.Id)
	assert.Zero(t, bad1.Id)
	assert.Zero(t, bad2.Id)
	assert.NotContains(t, es.EventMap, "https://opry.com/event/bad1")
	assert.NotContains(t, es.EventMap, "https://opry.com/event/bad2")

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM events WHERE link = ?", "https://opry.com/event/good").Scan(&count))
	assert.Equal(t, 1, count)
}

func TestLoadFromDB_SkipsBadTime(t *testing.T) {
	db := setupTestDB(t)
	now := time.Now().UTC().Truncate(time.Second)

	_, err := db.Exec(`INSERT INTO events (link, title, time, no_of_performers) VALUES (?, ?, ?, ?)`,
		"https://opry.com/event/1", "Good Event", now.Format(time.RFC3339), 0)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO events (link, title, time, no_of_performers) VALUES (?, ?, ?, ?)`,
		"https://opry.com/event/2", "Bad Time Event", "not-a-time", 0)
	require.NoError(t, err)

	es := NewEventStore()
	err = es.LoadFromDB(db)
	require.NoError(t, err)

	assert.Len(t, es.EventMap, 1)
	assert.Contains(t, es.EventMap, "https://opry.com/event/1")
}
