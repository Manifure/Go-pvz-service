package handler

import (
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/model"
	"Go-pvz-service/internal/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateItemRequest struct {
	PvzID string `json:"pvz_id"`
	Type  string `json:"type"`
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	// Проверяю наличие открытой приёмки
	acceptanceID, err := storage.GetOpenAcceptanceID(db.DB, req.PvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"message":"no open acceptance found"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"message":"failed to check open acceptance"}`, http.StatusInternalServerError)
		return
	}

	// Добавляю товар
	item := model.Item{
		ID:         uuid.New().String(),
		Type:       req.Type,
		ReceivedAt: time.Now(),
	}

	err = storage.AddItemToAcceptance(db.DB, item, acceptanceID)
	if err != nil {
		http.Error(w, `{"message":"failed to add item to acceptance"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
