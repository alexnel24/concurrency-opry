package server

import(
	"net/http"

	"github.com/alexnel24/concurrency-opry/internal/handlers"
)

func Routes(mux *http.ServeMux, h *handlers.Handler) http.Handler {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/scrape", h.Scrape)
	mux.HandleFunc("/update-db", h.UpdateDB)
	mux.HandleFunc("/mark-past-events", h.MarkPastEvents)
	mux.HandleFunc("/artist-performances", h.ArtistPerformances)
	mux.HandleFunc("/sessions", h.Sessions)

	return mux
}