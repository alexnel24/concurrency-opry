package store

import (
	"database/sql"
	"fmt"
	"sync"

	"OpryScrape/internal/models"
)

type PerformanceStore struct {
	mu        			sync.Mutex
	performanceMap  	map[string]*models.Performance
	newPerformancesCh	chan *models.Performance
}

func NewPerformanceStore() *PerformanceStore {
	return &PerformanceStore{
		performanceMap: make(map[string]*models.Performance),
		newPerformancesCh: make(chan *models.Performance),
	}
}

func (ps *PerformanceStore) AddPerformance(artistName string, event *models.Event) *models.Performance {
	
	ps.mu.Lock()
	defer ps.mu.Unlock()

	comboStr := fmt.Sprintf("%s-%s", artistName, event.Link)

	if Performance, exists := ps.performanceMap[comboStr]; exists {
		return Performance
	}

	performance := models.NewPerformance(artistName, event.Link)
	ps.performanceMap[performance.ComboString] = performance

	ps.newPerformancesCh <- performance
	
	event.AddOnePerformer()
	
	return performance
}

const performanceQuery = `
        SELECT id, event_link, artist_name, combo_string
        FROM performances;
    `
func (ps *PerformanceStore) LoadFromDB(db *sql.DB) error {
    rows, err := db.Query(performanceQuery)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var p models.Performance
        if err := rows.Scan(&p.Id, &p.EventLink, &p.ArtistName, &p.ComboString); err != nil {
            return err
        }

        ps.performanceMap[p.ComboString] = &p
    }

    return nil
}

const performanceInsert = `
        INSERT INTO performances (event_link, artist_name, combo_string)
        VALUES (?, ?, ?)
    `
func (ps *PerformanceStore) InsertPerformancesToDb(db *sql.DB, newPerformances []*models.Performance) error {
	tx, err := db.Begin()
    if err != nil {
        return err
    }
	defer tx.Rollback()

    stmt, err := tx.Prepare(performanceInsert)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, p := range newPerformances {
        result, err := stmt.Exec(p.EventLink, p.ArtistName, p.ComboString)
        if err != nil {
            fmt.Println("Error on Performance combo: ", p.ComboString)
            ps.mu.Lock()
            delete(ps.performanceMap, p.ComboString)
            ps.mu.Unlock()
            continue
        }

		id, err := result.LastInsertId()
		if err != nil {
            fmt.Println("Error getting ID for Performance combo: ", p.ComboString)
            ps.mu.Lock()
            delete(ps.performanceMap, p.ComboString)
            ps.mu.Unlock()
            continue
        }

		p.Id = id
    }

    return tx.Commit()
}