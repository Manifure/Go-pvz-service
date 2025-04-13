package handler

import (
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/storage"
	"encoding/json"
	"net/http"
)

func CreatePVZHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		City string `json:"city"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	pvz, err := storage.CreatePVZ(db.DB, input.City)
	if err != nil {
		if err == storage.ErrCityNotAllowed {
			http.Error(w, `{"message":"city not allowed"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"message":"failed to create pvz"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(pvz)
	if err != nil {
		http.Error(w, `{"message":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
