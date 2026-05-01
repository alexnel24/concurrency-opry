package store

import (
	"database/sql"
	"testing"

	"OpryScrape/internal/db/schema"
	"OpryScrape/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestEventStore_LoadFromDB_RoundTrip(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(schema.EventsSchema)
	require.NoError(t, err)

	store := NewEventStore()
	event := models.NewEvent("Opry 100", "https://opry.com/event/opry-100")

	err = store.InsertEventsToDb(db, []*models.Event{event})
	require.NoError(t, err)
	assert.NotZero(t, event.Id)

	freshStore := NewEventStore()
	err = freshStore.LoadFromDB(db)
	require.NoError(t, err)

	assert.Len(t, freshStore.EventMap, 1)
	loaded, ok := freshStore.EventMap["https://opry.com/event/opry-100"]
	require.True(t, ok)
	assert.Equal(t, "Opry 100", loaded.Title)
	assert.Equal(t, event.Id, loaded.Id)
}
