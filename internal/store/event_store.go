package store

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/alexnel24/concurrency-opry/internal/models"
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
        SELECT id, link, title, time, no_of_performers, upcoming
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
        if err := rows.Scan(&e.Id, &e.Link, &e.Title, &timeStr, &e.NoOfPerformers, &e.Upcoming); err != nil {
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
        INSERT INTO events (link, title, time, no_of_performers, upcoming)
        VALUES (?, ?, ?, ?, ?)
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
        result, err := stmt.Exec(e.Link, e.Title, e.Time, e.NoOfPerformers, e.Upcoming)
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

func (es *EventStore) UpdatePastEventsInDb(db *sql.DB) error {
	func() {
		es.mu.Lock()
		defer es.mu.Unlock()
		for _, e := range es.EventMap {
			if !e.Time.IsZero() && e.Time.Before(time.Now()) {
				e.Upcoming = false
			}
		}
	}()

	_, err := db.Exec(`UPDATE events SET upcoming = 0
		WHERE upcoming = 1
		AND time != '0001-01-01T00:00:00Z'
		AND time < datetime('now')`)
	return err
}

func (es *EventStore) SyncNoOfPerformersToDb(db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE events
		SET no_of_performers = (
			SELECT COUNT(*) FROM performances WHERE event_link = events.link
		)
		WHERE upcoming = 1
	`)
	return err
}

func (es *EventStore) SyncEventTimesToDb(db *sql.DB) error {
	candidates := func() []*models.Event {
		es.mu.Lock()
		defer es.mu.Unlock()
		c := make([]*models.Event, 0)
		for _, e := range es.EventMap {
			if e.Id != 0 && !e.Time.IsZero() {
				c = append(c, e)
			}
		}
		return c
	}()

	stmt, err := db.Prepare(`UPDATE events SET time = ? WHERE id = ? AND time = '0001-01-01T00:00:00Z'`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range candidates {
		if _, err := stmt.Exec(e.Time, e.Id); err != nil {
			fmt.Printf("Error syncing time for event id=%d: %s\n", e.Id, err.Error())
		}
	}
	return nil
}
