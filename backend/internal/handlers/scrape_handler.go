package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) Scrape(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.scraper.ScrapeOpry(ctx, h.stores)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error Scraping Opry Website: %s", err.Error())
	}

	h.stores.FlushAllOutstandingToDb()
	//ToDo add error handling for DB

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Scraping Complete")
}

//ToDo add helper function that takes in error and generates 500 response (used for both scraping and db)
