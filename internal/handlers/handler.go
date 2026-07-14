package handlers

import (
	"github.com/alexnel24/concurrency-opry/internal/services/performancefinder"
	"github.com/alexnel24/concurrency-opry/internal/services/scraping"
	"github.com/alexnel24/concurrency-opry/internal/services/watchlist"
	"github.com/alexnel24/concurrency-opry/internal/session"
	"github.com/alexnel24/concurrency-opry/internal/store"
)

type Handler struct {
	scraper           *scraping.Scraper
	watchlist         *watchlist.Watchlist
	performanceFinder *performancefinder.PerformanceFinder
	stores            *store.Stores
	sessionManager    *session.SessionManager
}

func New(scraper *scraping.Scraper, watchlist *watchlist.Watchlist, performanceFinder *performancefinder.PerformanceFinder, stores *store.Stores, sessionManager *session.SessionManager) *Handler {
	return &Handler{
		scraper:           scraper,
		watchlist:         watchlist,
		performanceFinder: performanceFinder,
		stores:            stores,
		sessionManager:    sessionManager,
	}
}