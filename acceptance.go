package handler

import (
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/model"
	"Go-pvz-service/internal/storage"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AcceptanceRequest struct {
	PvzID string        `json:"pvz_id"`
	Items []ItemRequest `json:"items"`
}

type ItemRequest struct {
	Type string `json:"type"`
}

func CreateAcceptanceHandler(w http.ResponseWriter, r *http.Request) {
	var req AcceptanceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	hasInProgress, err := storage.HasOpenAcceptance(db.DB, req.PvzID)
	if err != nil {
		http.Error(w, `{"message":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if hasInProgress {
		http.Error(w, `{"message":"previous acceptance is not closed"}`, http.StatusConflict)
		return
	}

	acceptance := model.Acceptance{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		PVZID:     req.PvzID,
		Status:    "in_progress",
	}

	for _, itemReq := range req.Items {
		acceptance.Items = append(acceptance.Items, model.Item{
			ID:         uuid.New().String(),
			ReceivedAt: acceptance.CreatedAt,
			Type:       itemReq.Type,
		})
	}

	err = storage.CreateAcceptanceWithItems(db.DB, &acceptance)
	if err != nil {
		http.Error(w, `{"message":"failed to create acceptance"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(acceptance); err != nil {
		http.Error(w, `{"message":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
