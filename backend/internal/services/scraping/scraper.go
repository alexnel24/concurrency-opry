package scraping

import (
	"context"
	"fmt"

	// "github.com/alexnel24/concurrency-opry/internal/models"
	"github.com/alexnel24/concurrency-opry/internal/store"
)

type Scraper struct {
	stores *store.Stores
}

func NewScraper(stores *store.Stores) *Scraper {
	return &Scraper{
		stores: stores,
	}
}

func (s *Scraper) ScrapeOpry(ctx context.Context, stores *store.Stores) (error) {
	months, _ := ScrapeAndGenerateMonths()
	fmt.Println("Done Scraping Months")

	//EXAMPLE OF FORCING A MONTH TO LIMIT EVENT SLICE FUNCTION BELOW
	//USED WHILE BUILDING OUT THE EVENT SLICE
	// dec, _ := models.NewMonth("December 2025")
	// tempMonthSlice := []models.Month{*dec}
	// s.ScrapeEvents(ctx, tempMonthSlice)

	s.ScrapeEvents(ctx, months)
	fmt.Println("Done Scraping Events")

	s.ScrapeArtistsAndPerformances(ctx)
	fmt.Println("Done Scraping Artists and Performances")

	s.stores.FlushAllOutstandingToDb()
	if err := s.stores.SyncEventTimesToDb(); err != nil {
		fmt.Println("Error syncing event times:", err)
	}

	return nil
}