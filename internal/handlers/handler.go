package handlers

import (
	"OpryScrape/internal/services/scraping"
	"OpryScrape/internal/store"
)

type Handler struct {
	scraper *scraping.Scraper
	stores  *store.Stores
}

func New(scraper *scraping.Scraper, stores *store.Stores) *Handler {
	return &Handler{
		scraper: scraper,
		stores: stores,
	}
}