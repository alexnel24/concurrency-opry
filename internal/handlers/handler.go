package handlers

import (
	"github.com/alexnel24/concurrency-opry/internal/services/scraping"
	"github.com/alexnel24/concurrency-opry/internal/store"
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