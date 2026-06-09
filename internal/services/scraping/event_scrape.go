package scraping

import (
	"context"
	"fmt"

	"github.com/alexnel24/concurrency-opry/internal/models"
	"github.com/alexnel24/concurrency-opry/internal/parse"

	"github.com/gocolly/colly"
	"golang.org/x/sync/errgroup"
)

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
					link := e.Attr("href")
					s.stores.EventStore.AddEvent(title, link, parse.ParseTimeFromLink(link))
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
