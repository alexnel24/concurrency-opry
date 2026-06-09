package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}