package store

import (
	"github.com/alexnel24/concurrency-opry/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Stores struct {
	db               *sql.DB
	EventStore       *EventStore
	ArtistStore      *ArtistStore
	PerformanceStore *PerformanceStore
	FlushToDbCh		 chan struct{}
}

func InitStores(db *sql.DB) *Stores {
	store := &Stores{
		db:               db,
		EventStore:       NewEventStore(),
		ArtistStore:      NewArtistStore(),
		PerformanceStore: NewPerformanceStore(),
		FlushToDbCh:      make(chan struct{}),	
	}
	store.loadFromDB()
	return store
}

//ToDo: Ideally return bool and err or just err
func (s *Stores) loadFromDB() {
	s.EventStore.LoadFromDB(s.db)
	s.ArtistStore.LoadFromDB(s.db)
	s.PerformanceStore.LoadFromDB(s.db)
}

func (s *Stores) StartBackgroundDBWorker(ctx context.Context, batchSize, flushEverySeconds int) {
	
	go func() {
		ticker := time.NewTicker(time.Duration(flushEverySeconds) * time.Second)
		defer ticker.Stop()

		artistBatch := make([]*models.Artist, 0 ,batchSize)
		eventBatch := make([]*models.Event, 0, batchSize)
		performanceBatch := make([]*models.Performance, 0, batchSize)

		//setup funcs
		flushArtists := func(){
			if len(artistBatch) == 0 {return}
			err := s.ArtistStore.InsertArtistsToDb(s.db, artistBatch)
			if err != nil {
				fmt.Println("Error pushing artists: ", err.Error())
			} else {
				fmt.Println("Pushed artists to db")
			}
			artistBatch = artistBatch[:0]
		}

		flushEvents := func(){
			if len(eventBatch) == 0 {return}
			err := s.EventStore.InsertEventsToDb(s.db, eventBatch)
			if err != nil {
				fmt.Println("Error pushing events: ", err.Error())
				fmt.Println("")
			} else {
				fmt.Println("Pushed events to db")
			}

			eventBatch = eventBatch[:0]
		}

		flushPerformances := func(){
			if len(performanceBatch) == 0 {return}
			err := s.PerformanceStore.InsertPerformancesToDb(s.db, performanceBatch)
			if err != nil {
				fmt.Println("Error pushing performances: ", err.Error())
			} else {
				fmt.Println("Pushed performances to db")
			}
			performanceBatch = performanceBatch[:0]
		}

		flushAll := func() {
			flushArtists()
			flushEvents()
			flushPerformances()
		}
		//end of setting up functions

		//listen for channel activity
		for {
			select{
			
			case artist, ok := <-s.ArtistStore.newArtistsCh:
				if !ok {
					flushArtists()
					s.ArtistStore.newArtistsCh = nil

					break 
				}

				artistBatch = append(artistBatch, artist)
				if len(artistBatch) >= batchSize {
					flushArtists()
				}

			case event, ok := <-s.EventStore.newEventsCh:
				if !ok {
					flushEvents()
					s.EventStore.newEventsCh = nil

					break
				}

				eventBatch = append(eventBatch, event)
				if len(eventBatch) >= batchSize {
					flushEvents()
				}
			
			case performance, ok := <-s.PerformanceStore.newPerformancesCh:
				if !ok {
					flushPerformances()
					s.PerformanceStore.newPerformancesCh = nil

					break
				}

				performanceBatch = append(performanceBatch, performance)
				if len(performanceBatch) >= batchSize {
					flushPerformances()
				}

			case <-ticker.C:
				flushAll()

			case <-s.FlushToDbCh:
				flushAll()

			case <-ctx.Done():
				flushAll()

				return
			}
		}
	}()
}

func (s *Stores) FlushAllOutstandingToDb(){
	s.FlushToDbCh <- struct{}{}
}

func (s *Stores) DB() *sql.DB {
	return s.db
}

func (s *Stores) SyncEventTimesToDb() error {
	return s.EventStore.SyncEventTimesToDb(s.db)
}

func (s *Stores) SyncNoOfPerformersToDb() error {
	return s.EventStore.SyncNoOfPerformersToDb(s.db)
}

