package handler

import (
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

type CloseAcceptanceRequest struct {
	PvzID string `json:"pvz_id"`
}

func CloseAcceptanceHandler(w http.ResponseWriter, r *http.Request) {
	var req CloseAcceptanceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PvzID == "" {
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
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

	err = storage.CloseAcceptance(db.DB, acceptanceID)
	if err != nil {
		http.Error(w, `{"message":"failed to close acceptance"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"acceptance closed successfully"}`))
}
