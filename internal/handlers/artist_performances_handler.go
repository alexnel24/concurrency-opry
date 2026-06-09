package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type artistPerformanceResponse struct {
	ArtistName string `json:"artist_name"`
	EventLink  string `json:"event_link"`
	EventTitle string `json:"event_title"`
	EventTime  string `json:"event_time"`
	Upcoming   bool   `json:"upcoming"`
}

func (h *Handler) ArtistPerformances(w http.ResponseWriter, r *http.Request) {
	artists := r.URL.Query()["artist"]
	if len(artists) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "missing required query param: artist")
		return
	}

	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "all"
	}
	switch filter {
	case "all", "upcoming", "past":
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid filter value: must be one of all, upcoming, past")
		return
	}

	performances := h.stores.PerformanceStore.GetAllByArtists(artists)

	resp := make([]artistPerformanceResponse, 0, len(performances))
	for _, p := range performances {
		event, ok := h.stores.EventStore.EventMap[p.EventLink]
		if !ok {
			continue
		}
		switch filter {
		case "upcoming":
			if !event.Upcoming {
				continue
			}
		case "past":
			if event.Upcoming {
				continue
			}
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
		fmt.Println("error encoding artist performances response: ", err)
	}
}
