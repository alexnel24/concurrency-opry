package handlers

import (
	"github.com/alexnel24/concurrency-opry/internal/services/scraping"
	"github.com/alexnel24/concurrency-opry/internal/session"
	"github.com/alexnel24/concurrency-opry/internal/store"
)

type Handler struct {
	scraper        *scraping.Scraper
	stores         *store.Stores
	sessionManager *session.SessionManager
}

func New(scraper *scraping.Scraper, stores *store.Stores, sessionManager *session.SessionManager) *Handler {
	return &Handler{
		scraper:        scraper,
		stores:         stores,
		sessionManager: sessionManager,
	}
}