package scraping

import (
	"context"
	"fmt"
	"time"

	"github.com/alexnel24/concurrency-opry/internal/parse"

	"github.com/gocolly/colly"
	"golang.org/x/sync/errgroup"
)

func (s *Scraper) ScrapeArtistsAndPerformances(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	baseCollector := colly.NewCollector(
		colly.AllowedDomains("opry.com", "www.opry.com"),
	)

	baseCollector.SetRequestTimeout(30 * time.Second)

	for _, event := range s.stores.EventStore.EventMap {
		event := event

		g.Go(func() error {
			c := baseCollector.Clone()

			c.OnHTML("div.artist_list h3.title span", func(e *colly.HTMLElement) {
				artist := s.stores.ArtistStore.AddArtist(e.Text)
				s.stores.PerformanceStore.AddPerformance(artist.Name, event)
			})

			if event.Time.IsZero() {
				c.OnHTML("li", func(e *colly.HTMLElement) {
					t := parse.ParseDateTimeFromText(e.Text)
					if !t.IsZero() {
						s.stores.EventStore.UpdateEventTime(event.Link, t)
					}
				})
			}

			err := c.Visit(event.Link)
			if err != nil {
				fmt.Println("Error encountered while looking for artists and performances: ", err.Error())
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err // first non-nil error
	}

	return nil
}
