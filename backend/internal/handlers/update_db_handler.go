package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) UpdateDB(w http.ResponseWriter, r *http.Request){
	h.stores.FlushAllOutstandingToDb()
	//ToDo need to check if the write to DB was successful rather than just saying it was
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success Updating the DB")
}