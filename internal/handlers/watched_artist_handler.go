package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexnel24/concurrency-opry/internal/util/ownerid"
)

func writeBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "%s", msg)
}

func (h *Handler) WatchedArtists(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listWatchedArtists(w, r)
	case http.MethodPost:
		h.addWatchedArtist(w, r)
	case http.MethodDelete:
		h.removeWatchedArtist(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listWatchedArtists(w http.ResponseWriter, r *http.Request) {
	ownerId, err := ownerid.Parse(r.URL.Query().Get("ownerId"))
	if err != nil {
		writeBadRequest(w, "invalid ownerId value: must be an integer")
		return
	}

	watched, err := h.watchlist.ListWatchedArtists(ownerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error retrieving watched artists")
		return
	}

	resp := make([]watchedArtistResponse, 0, len(watched))
	for _, wa := range watched {
		resp = append(resp, watchedArtistResponse{
			Id:       wa.Id,
			Name:     wa.Name,
			ArtistId: wa.ArtistId,
			OwnerId:  wa.OwnerId,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println("error encoding watched artists response: ", err)
	}
}

func (h *Handler) addWatchedArtist(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeBadRequest(w, "missing required query param: name")
		return
	}

	ownerId, err := ownerid.Parse(r.URL.Query().Get("ownerId"))
	if err != nil {
		writeBadRequest(w, "invalid ownerId value: must be an integer")
		return
	}

	result, err := h.watchlist.AddWatchedArtist(name, ownerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error adding watched artist")
		return
	}

	resp := watchedArtistResponse{
		Id:            result.Entry.Id,
		Name:          result.Entry.Name,
		ArtistId:      result.Entry.ArtistId,
		OwnerId:       result.Entry.OwnerId,
		AlreadyExists: result.AlreadyExists,
	}

	// 200 if it already existed
	if result.AlreadyExists {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "artist %s already on watchlist", result.Entry.Name)
		return
	}

	// 201 if newly created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println("error encoding watched artists response: ", err)
	}
}

func (h *Handler) removeWatchedArtist(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeBadRequest(w, "missing required query param: name")
		return
	}

	ownerId, err := ownerid.Parse(r.URL.Query().Get("ownerId"))
	if err != nil {
		writeBadRequest(w, "invalid ownerId value: must be an integer")
		return
	}

	found, err := h.watchlist.RemoveWatchedArtist(name, ownerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error removing watched artist")
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
