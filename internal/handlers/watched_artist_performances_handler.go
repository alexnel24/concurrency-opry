package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexnel24/concurrency-opry/internal/models"
	"github.com/alexnel24/concurrency-opry/internal/util/ownerid"
)

func (h *Handler) WatchedArtistPerformances(w http.ResponseWriter, r *http.Request) {
	ownerId, err := ownerid.Parse(r.URL.Query().Get("ownerId"))
	if err != nil {
		writeBadRequest(w, "invalid ownerId value: must be an integer")
		return
	}

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = models.FilterAll
	}
	switch filter {
	case models.FilterAll, models.FilterUpcoming, models.FilterPast:
	default:
		writeBadRequest(w, fmt.Sprintf("invalid filter value: must be one of %s, %s, %s", models.FilterAll, models.FilterUpcoming, models.FilterPast))
		return
	}

	performances, err := h.watchlist.GetWatchedArtistPerformances(ownerId, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error retrieving watched artist performances")
		return
	}

	resp := make([]artistPerformanceResponse, 0, len(performances))
	for _, p := range performances {
		event, ok := h.stores.EventStore.EventMap[p.EventLink]
		if !ok {
			continue
		}
		resp = append(resp, artistPerformanceResponse{
			ArtistName: p.ArtistName,
			EventLink:  p.EventLink,
			EventTitle: event.Title,
			EventTime:  event.Time.Format(time.RFC3339),
			Upcoming:   event.Upcoming,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println("error encoding watched artist performances response: ", err)
	}
}
