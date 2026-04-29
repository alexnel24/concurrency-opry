package scraping

import (
	"context"
	"fmt"

	"OpryScrape/internal/models"

	"github.com/gocolly/colly"
	"golang.org/x/sync/errgroup"
)

/*
Example comment - for Rockbot in case they want to see I know how to do it
Scrape loops through the given months currently on the opry website, and caputres new events (shows)
Parameters:
	ctx: The context to be used for cancelling (REST API)
	months: the months to search for new events
Returns: 
	ERROR only as the new events are stored within the eventStore
*/
func (s *Scraper) ScrapeEvents(ctx context.Context, months []models.Month) error {
	g, _ := errgroup.WithContext(ctx)

	baseCollector := colly.NewCollector(
		colly.AllowedDomains("opry.com", "www.opry.com"),
	)

	for _, month := range months {
		month := month

		g.Go(func() error {
			select {

			case <- ctx.Done():
				fmt.Println("ScrapeEvents cancelled during month: ", month)
				return ctx.Err()
			
			default:
				c := baseCollector.Clone()

				c.OnHTML(".eventList__wrapper.list h3.title a", func(e *colly.HTMLElement) {
					title := e.Text
					s.stores.EventStore.AddEvent(title, e.Attr("href"))
				})

				err := c.Visit("https://www.opry.com/events/filtered/" + month.Year + "/" + month.MonthStr)
				if err != nil {
					//ToDo: error handling (skip month)
					fmt.Println("error visting MONTH specific page")
					return err
				}
				return nil
			}
		})
	}

	err := g.Wait()
	if err != nil {
		return err // first non-nil error
	}

	return nil
}
