package handler

import (
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

type DeleteItemRequest struct {
	PvzID string `json:"pvz_id"`
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	acceptanceID, err := storage.GetOpenAcceptanceID(db.DB, req.PvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"message":"no open acceptance found"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"message":"failed to check open acceptance"}`, http.StatusInternalServerError)
		return
	}

	err = storage.DeleteLastItemFromAcceptance(db.DB, acceptanceID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"message":"no items found in acceptance"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"message":"failed to delete item"}`, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"item deleted"}`))
}
