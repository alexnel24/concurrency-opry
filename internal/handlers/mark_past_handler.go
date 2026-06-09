package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) MarkPastEvents(w http.ResponseWriter, r *http.Request) {
	err := h.stores.EventStore.UpdatePastEventsInDb(h.stores.DB())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marking past events: %s", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Past events marked successfully")
}
