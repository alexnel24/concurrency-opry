package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type createSessionResponse struct {
	SessionId string `json:"session_id"`
}

func (h *Handler) Sessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createSession(w, r)
	case http.MethodDelete:
		h.deleteSession(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	id := h.sessionManager.Create()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createSessionResponse{SessionId: id}); err != nil {
		fmt.Println("error encoding create session response: ", err)
	}
}

func (h *Handler) deleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "missing X-Session-ID header")
		return
	}

	if !h.sessionManager.Delete(sessionID) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "session not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
