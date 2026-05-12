package scraping

import (
	"errors"
	"fmt"

	"github.com/alexnel24/concurrency-opry/internal/models"

	"github.com/gocolly/colly"
)

func ScrapeAndGenerateMonths() ([]models.Month, error){
	c := colly.NewCollector(
		colly.AllowedDomains("opry.com", "www.opry.com"),
	)

	var monthSlice []models.Month

	c.OnHTML("a.event_filter_item", func(e *colly.HTMLElement) {
		monthObject, err := models.NewMonth(e.Text)
		if err != nil {
			//ToDo: figure out what opry might send in as bad month, and how to handle
			fmt.Println("Error encountered: ", err.Error())
		}
		monthSlice = append(monthSlice, *monthObject)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit("https://www.opry.com/events")
	if err != nil {
		return nil, errors.New("Error visiting first website: " + err.Error())
	}

	return monthSlice, nil
}