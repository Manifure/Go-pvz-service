package handler

import (
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/model"
	"Go-pvz-service/internal/storage"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetPVZDataHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	params := model.PVZQueryParams{
		From:   from,
		To:     to,
		Limit:  limit,
		Offset: offset,
	}

	pvzList, err := storage.GetPVZWithAcceptancesFiltered(db.DB, params)
	if err != nil {
		http.Error(w, `{"message":"failed to fetch pvz data"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pvzList); err != nil {
		http.Error(w, `{"message":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
