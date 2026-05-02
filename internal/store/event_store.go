package store

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"OpryScrape/internal/models"
)

type EventStore struct {
	mu        		sync.Mutex
	EventMap  		map[string]*models.Event
	newEventsCh		chan *models.Event
}

func NewEventStore() *EventStore {
	return &EventStore{
		EventMap: make(map[string]*models.Event),
		newEventsCh: make(chan *models.Event, 100),
	}
}

func (es *EventStore) AddEvent(title, link string, t time.Time) *models.Event {
	es.mu.Lock()
	defer es.mu.Unlock()

	if event, exists := es.EventMap[link]; exists {
		return event
	}

	event := models.NewEvent(title, link, t)
	es.EventMap[link] = event
	es.newEventsCh <- event

	return event
}

func (es *EventStore) UpdateEventTime(link string, t time.Time) {
	es.mu.Lock()
	defer es.mu.Unlock()
	if event, exists := es.EventMap[link]; exists && event.Time.IsZero() {
		event.Time = t
	}
}

const eventQuery = `
        SELECT id, link, title, time, no_of_performers
        FROM events;
    `
func (es *EventStore) LoadFromDB(db *sql.DB) error {
    rows, err := db.Query(eventQuery)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var e models.Event
		var timeStr string
        if err := rows.Scan(&e.Id, &e.Link, &e.Title, &timeStr, &e.NoOfPerformers); err != nil {
            return err
        }

        e.Time, err = time.Parse(time.RFC3339, timeStr)
        if err != nil {
            fmt.Println("unable to parse time on event: ", e.Id)
            continue
        }

        es.EventMap[e.Link] = &e
    }

    return nil
}


const eventInsert = `
        INSERT INTO events (link, title, time, no_of_performers)
        VALUES (?, ?, ?, ?)
    `
func (es *EventStore) InsertEventsToDb(db *sql.DB, newEvents []*models.Event) error {
	tx, err := db.Begin()
    if err != nil {
        return err
    }
	defer tx.Rollback()

    stmt, err := tx.Prepare(eventInsert)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, e := range newEvents {
        result, err := stmt.Exec(e.Link, e.Title, e.Time, e.NoOfPerformers)
        if err != nil {
            fmt.Println("Error on Event link: ", e.Link)
            es.mu.Lock()
            delete(es.EventMap, e.Link)
            es.mu.Unlock()
            continue
        }

        id, err := result.LastInsertId()
        if err != nil {
            fmt.Println("Error getting ID for Event link: ", e.Link)
            es.mu.Lock()
            delete(es.EventMap, e.Link)
            es.mu.Unlock()
            continue
        }

        e.Id = id
    }

    return tx.Commit()
}
